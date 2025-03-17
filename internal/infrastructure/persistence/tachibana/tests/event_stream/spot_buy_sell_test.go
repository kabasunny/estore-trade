// internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_sell_test.go
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

func TestEventStreamSpotBuyAndSell(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (現物成行買い -> 現物売り)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// Dispatcher の作成
		dispatcher := dispatcher.NewOrderEventDispatcher(client.GetLogger())
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), dispatcher)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// --- 現物成行買い注文 ---
		buyMarketCode := "00" // 買い注文の市場コードを保持
		buyOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: buyMarketCode, // 東証
		}

		// 買い注文用のチャネルと購読ID
		buyEventCh := make(chan *domain.OrderEvent, 10)
		buyOrderID := "buyOrder"
		dispatcher.Subscribe(buyOrderID, buyEventCh)
		defer dispatcher.Unsubscribe(buyOrderID, buyEventCh)

		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		// 購読IDと注文IDを関連づける
		dispatcher.RegisterOrderID(placedBuyOrder.TachibanaOrderID, buyOrderID)
		// --- 買い注文の約定確認 ---
		var buyEvent *domain.OrderEvent
		timeout := time.After(60 * time.Second)
	BuyOrderLoop: // ラベルを追加
		for {
			select {
			case event := <-buyEventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedBuyOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" { //全部約定
						t.Logf("Buy order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						buyEvent = event   // 約定イベントを保存
						break BuyOrderLoop // 買い注文の約定確認ループを抜ける
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Buy order failed. Status: %s", event.Order.Status)
						return
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Buy order execution event not received")
				return //追加
			}
		}

		// --- 現物売り注文 (返済) ---
		sellOrder := &domain.Order{
			Symbol:     buyEvent.Order.Symbol,         // 買い注文と同じ銘柄
			Side:       "short",                       // 売り
			OrderType:  "market",                      // 成行 (または指値)
			Quantity:   buyEvent.Order.FilledQuantity, // 買い注文の約定数量
			MarketCode: buyMarketCode,                 // 買い注文時に指定した市場コード
		}

		//売り注文用のチャネルと購読ID
		sellEventCh := make(chan *domain.OrderEvent, 10)
		sellOrderID := "sellOrder"
		dispatcher.Subscribe(sellOrderID, sellEventCh)
		defer dispatcher.Unsubscribe(sellOrderID, sellEventCh)

		placedSellOrder, err := client.PlaceOrder(ctx, sellOrder)
		if err != nil {
			t.Fatalf("Failed to place sell order: %v", err)
		}
		assert.NotNil(t, placedSellOrder)

		//購読IDと注文IDを関連付ける
		dispatcher.RegisterOrderID(placedSellOrder.TachibanaOrderID, sellOrderID)
		// --- 売り注文の約定確認 ---
		timeout = time.After(60 * time.Second) // タイムアウトをリセット
		for {
			select {
			case event := <-sellEventCh:
				if event == nil { //nilチェック
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedSellOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" { //注文失敗
						t.Logf("Sell order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						return // テスト成功 (売り注文の約定を確認)
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Sell order failed. Status: %s", event.Order.Status)
						return
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Sell order execution event not received")
				return //追加
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_sell_test.go

// 立花証券APIにログイン。
// イベントストリームを開始。
// 現物株の成行買い注文を発注。
// 買い注文の約定イベントを待ち、受信したら内容を確認。
// 買い注文の約定情報を使って、現物売り注文 (返済注文) を発注。
// 売り注文の約定イベントを待ち、受信したら内容を確認。
// すべての処理が成功したら、テスト成功。
// タイムアウトが発生したり、エラーが発生したりした場合は、テスト失敗。
