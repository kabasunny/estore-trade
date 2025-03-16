package tachibana_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamSpotBuy(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (現物成行買い)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		dispatcher := dispatcher.NewOrderEventDispatcher(client.GetLogger())
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), dispatcher)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		time.Sleep(1 * time.Second) // イベントストリーム接続確立を待つ

		eventCh := make(chan *domain.OrderEvent, 10)
		subscriberID := "test-subscriber"

		dispatcher.Subscribe(subscriberID, eventCh)
		defer dispatcher.Unsubscribe(subscriberID, eventCh)

		order := &domain.Order{
			Symbol:     "7974",
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00",
		}
		placedOrder, err := client.PlaceOrder(ctx, order)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err)
		}
		assert.NotNil(t, placedOrder)
		placedOrderID := placedOrder.TachibanaOrderID

		// 購読IDを更新 (Unsubscribe/Subscribe)
		dispatcher.Unsubscribe(subscriberID, eventCh)
		dispatcher.Subscribe(placedOrderID, eventCh)
		defer dispatcher.Unsubscribe(placedOrderID, eventCh)

		timeout := time.After(30 * time.Second) // タイムアウト時間を長めに設定

		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				t.Logf("Received event: %+v", event) // t.Logf を使用

				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedOrderID {
					if event.Order.Status == "1" || event.Order.Status == "3" {
						if event.Order.FilledQuantity > 0 {
							t.Logf("Order partially or fully executed. Status: %s, Executed Quantity:%d", event.Order.Status, event.Order.FilledQuantity)
							return
						}
					} else if event.Order.Status == "4" || event.Order.Status == "5" {
						t.Fatalf("Order failed. Status: %s", event.Order.Status)
						return
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Execution event not received after 30 seconds") // タイムアウトメッセージを修正
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamSpotBuy
