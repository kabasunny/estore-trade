// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_buy_with_stop_combined_test.go
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

func TestEventStreamCreditBuyWithStopCombined(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (信用新規成行買い + 逆指値売り - 同時注文)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		eventCh := make(chan *domain.OrderEvent, 10)
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), eventCh)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		time.Sleep(3 * time.Second) // イベントストリーム接続確立を待つ
		// --- 信用新規成行買い + 逆指値売り注文 ---
		buyOrder := &domain.Order{
			Symbol:       "7974", // 例: 任天堂
			Side:         "long",
			OrderType:    "stop",        // 通常 + 逆指値
			TradeType:    "credit_open", // 信用新規
			Quantity:     100,
			MarketCode:   "00",   // 東証
			Price:        0,      // 成行
			TriggerPrice: 9000.0, // 例: 現在価格が10000円として、9000円以下になったら
		}
		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		buyOrderID := placedBuyOrder.TachibanaOrderID

		// --- 注文状態の確認（逆指値なので、最初は受付済） ---
		timeout := time.After(60 * time.Second)
		var orderStatus string
	OrderLoop:
		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == buyOrderID {
					orderStatus = event.Order.Status // ステータスを更新
					// 逆指値注文のステータス遷移: 1(受付済) -> [トリガー] -> 3(一部約定) or 2(全部約定)
					if orderStatus == "1" {
						t.Logf("Order status: %s (Waiting for trigger)", orderStatus)
						break OrderLoop // 受付確認後、ループを抜ける
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Buy order execution event not received.  Last status: %s", orderStatus)
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamCreditBuyWithStopCombined
