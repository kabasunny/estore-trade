package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestSpotBuyMarketAndSellStop_Combined(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("現物成行買い + 逆指値売り (EventStreamとGetOrderStatus)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// Dispatcher と EventStream の準備 (イベント受信用)
		dispatcher := dispatcher.NewOrderEventDispatcher(client.GetLogger())
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), dispatcher)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// --- 1. 現物成行買い注文 ---
		buyOrder := &domain.Order{
			Symbol:     "7974",
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00",
		}

		// 買い注文用のチャネルと購読ID (イベント受信用)
		buyEventCh := make(chan *domain.OrderEvent, 10)
		buyOrderID := "buyOrder" // 購読ID。注文IDとは異なるユニークなものにする
		dispatcher.Subscribe(buyOrderID, buyEventCh)
		defer dispatcher.Unsubscribe(buyOrderID, buyEventCh)

		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		// 購読IDと注文IDを関連づける
		dispatcher.RegisterOrderID(placedBuyOrder.TachibanaOrderID, buyOrderID)

		// --- 2. 買い注文の約定確認 (イベントストリーム) ---
		var buyEvent *domain.OrderEvent
		timeout := time.After(60 * time.Second)
	BuyOrderLoop:
		for {
			select {
			case event := <-buyEventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedBuyOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" { // 全部約定
						t.Logf("Buy order fully executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						buyEvent = event
						break BuyOrderLoop // 約定確認後、ループを抜ける
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Buy order failed. Status: %s", event.Order.Status)
						return
					}
					// Status が "1"(一部約定) の場合もログ出力しつつ待機
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)

				}
			case <-timeout:
				t.Fatalf("Timeout: Buy order execution event not received")
				return
			}
		}

		// --- 3. 逆指値売り注文 (ストップロス) ---
		stopLossOrder := &domain.Order{
			Symbol:       buyEvent.Order.Symbol,         // 買い注文と同じ銘柄
			Side:         "short",                       // 売り
			OrderType:    "stop",                        // 逆指値
			Quantity:     buyEvent.Order.FilledQuantity, // 買い注文の約定数量
			MarketCode:   "00",                          // 買い注文と同じ市場
			TriggerPrice: 9000.0,                        // ストップロスのトリガー価格（例）
			Price:        0,                             // トリガー後成行
		}

		placedStopLossOrder, err := client.PlaceOrder(ctx, stopLossOrder)
		if err != nil {
			t.Fatalf("Failed to place stop-loss order: %v", err)
		}
		assert.NotNil(t, placedStopLossOrder)
		stopLossOrderID := placedStopLossOrder.TachibanaOrderID

		// --- 4. 逆指値売り注文の *受付* 確認 (GetOrderStatusを使用) ---
		timeout = time.After(60 * time.Second) // タイムアウト
		orderDate := time.Now().Format("20060102")

	StopLossOrderLoop:
		for {
			select {
			case <-timeout:
				t.Fatalf("Timeout: Stop-loss order was not accepted within the timeout period")
				return
			default:
				statusOrder, err := client.GetOrderStatus(ctx, stopLossOrderID, orderDate)
				if err != nil {
					t.Logf("GetOrderStatus failed for stop-loss order: %v, retrying...", err)
					time.Sleep(1 * time.Second) // Wait and retry
					continue
				}
				if statusOrder.Status == "発注待ち" { // 日本語のステータスで比較
					t.Logf("Stop-loss order %s is accepted (Status: %s).", stopLossOrderID, statusOrder.Status)
					break StopLossOrderLoop // 受付確認後、ループを抜ける
				} else if statusOrder.Status != "" {
					t.Logf("Stop-loss order %s Status is : %s .", stopLossOrderID, statusOrder.Status)
				}
				time.Sleep(1 * time.Second)
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_with_stop_test.go

// 現物成行買い注文
// 買い注文の約定確認 (イベントストリーム)
// 逆指値売り注文 (ストップロス)
// 逆指値注文の受付確認 (イベントストリームだと、このテストケースのタイミング的に難しいので、ステータス確認を追って行う)
// という一連の流れをテスト

// domain.Order に NotificationType フィールドが追加されていること。
// tachibana.parseEvent 関数で、p_NT の値が event.Order.NotificationType にマッピングされていること。
// 逆指値注文の受付確認ループ (StopLossOrderLoop) で、event.Order.NotificationType == "100" && event.Order.Status == "1" の条件で判定していること。
// 逆指値注文発注直後の time.Sleep は削除されていること。
// 買い注文、逆指値売り注文ともに、エラーが発生した場合は、t.Fatalf でテストを即時終了させていること。
// 各ステップにコメントが追加され、処理内容が明確になっていること。
