// internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_with_stop_separate_test.go
package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamSpotBuyWithStopSeparate(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (現物成行買い + 逆指値売り - 分割注文)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		eventCh := make(chan *domain.OrderEvent, 100) // バッファを増やす
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		time.Sleep(3 * time.Second) // イベントストリーム接続確立を待つ

		// --- 1. 現物成行買い注文 ---
		buyOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00", // 東証
		}
		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		buyOrderID := placedBuyOrder.TachibanaOrderID

		// --- 2. 買い注文の約定確認 ---
		var buyEvent *domain.OrderEvent
		timeout := time.After(60 * time.Second)
	BuyOrderLoop:
		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == buyOrderID {
					if (event.Order.Status == "1" || event.Order.Status == "3") && event.Order.FilledQuantity > 0 {
						t.Logf("Buy order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						buyEvent = event
						break BuyOrderLoop // 約定確認後、ループを抜ける
					}
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
			TriggerPrice: 9000.0,                        // ストップロスのトリガー価格
			Price:        0,                             // トリガー後成行
		}
		placedStopLossOrder, err := client.PlaceOrder(ctx, stopLossOrder)
		if err != nil {
			t.Fatalf("Failed to place stop-loss order: %v", err)
		}
		assert.NotNil(t, placedStopLossOrder)
		stopLossOrderID := placedStopLossOrder.TachibanaOrderID

		// --- 4. 逆指値売り注文の受付確認 ---
		//time.Sleep(1 * time.Second) // 不要

		timeout = time.After(60 * time.Second) // タイムアウトをリセット
		var stopLossOrderStatus string
	StopLossOrderLoop:
		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == stopLossOrderID {
					stopLossOrderStatus = event.Order.Status
					// NotificationType が "100" (注文状態変更) かつ status が "1" (受付済) ならOK
					if event.Order.NotificationType == "100" && stopLossOrderStatus == "1" {
						t.Logf("Stop-loss order accepted. Status: %s, NotificationType: %s", stopLossOrderStatus, event.Order.NotificationType)
						break StopLossOrderLoop // 受付確認後、ループを抜ける
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Stop-loss order was not accepted. Last status: %s", stopLossOrderStatus)
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSpotBuyWithStopSeparate

// 現物成行買い注文
// 買い注文の約定確認 (イベントストリーム)
// 逆指値売り注文 (ストップロス)
// 逆指値注文の受付確認 (イベントストリーム)
// という一連の流れをテスト

// domain.Order に NotificationType フィールドが追加されていること。
// tachibana.parseEvent 関数で、p_NT の値が event.Order.NotificationType にマッピングされていること。
// 逆指値注文の受付確認ループ (StopLossOrderLoop) で、event.Order.NotificationType == "100" && event.Order.Status == "1" の条件で判定していること。
// 逆指値注文発注直後の time.Sleep は削除されていること。
// 買い注文、逆指値売り注文ともに、エラーが発生した場合は、t.Fatalf でテストを即時終了させていること。
// 各ステップにコメントが追加され、処理内容が明確になっていること。
