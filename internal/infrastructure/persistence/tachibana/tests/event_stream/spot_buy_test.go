// internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_test.go
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

		order := &domain.Order{
			Symbol:     "7974",
			Side:       "long",
			OrderType:  "market",
			Quantity:   100,
			MarketCode: "00",
		}

		// 注文発注 *前* に eventCh を作成し、"system" で購読を開始
		eventCh := make(chan *domain.OrderEvent, 10)
		subscriberID := "test-subscriber" // 任意の文字列でOK
		dispatcher.Subscribe(subscriberID, eventCh)
		defer dispatcher.Unsubscribe(subscriberID, eventCh) // テスト終了時に購読解除

		// 注文発注
		placedOrder, err := client.PlaceOrder(ctx, order)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err)
		}
		assert.NotNil(t, placedOrder)

		// TachibanaOrderID と subscriberID のマッピングを登録
		placedOrderID := placedOrder.TachibanaOrderID
		dispatcher.RegisterOrderID(placedOrderID, subscriberID) // ★ここが重要

		timeout := time.After(60 * time.Second) // タイムアウト時間を長めに設定(60秒)

		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				t.Logf("Received event: %+v", event) //詳細ログ

				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedOrderID {
					// p_EXST == "2" (全部約定) を確認。
					if event.Order.ExecutionStatus == "2" {
						t.Logf("Order fully executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						return // テスト成功
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Order failed. Status: %s", event.Order.Status)
						return
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)

				}
			case <-timeout:
				t.Fatalf("Timeout: Execution event not received after 60 seconds") // タイムアウト時間を修正
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/spot_buy_test.go
