// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_buy_test.go

package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher をインポート
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamCreditBuy(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (信用新規買い)", func(t *testing.T) {
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// dispatcher の作成
		dispatcher := dispatcher.NewOrderEventDispatcher(client.GetLogger())
		eventStream := tachibana.NewEventStream(client, client.GetConfig(), client.GetLogger(), dispatcher)

		go func() {
			err := eventStream.Start(ctx)
			if err != nil {
				t.Errorf("EventStream Start returned error: %v", err)
			}
		}()
		defer eventStream.Stop()

		// --- 信用新規買い注文 ---
		creditBuyOrder := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long",
			OrderType:  "market", // または "limit"
			Quantity:   100,
			MarketCode: "00",          // 東証
			TradeType:  "credit_open", // 信用新規
			// Price:      13000,      // 指値の場合 (例)
		}

		// 買い注文用のチャネルと購読ID
		buyEventCh := make(chan *domain.OrderEvent, 10)
		buyOrderID := "buyOrder"
		dispatcher.Subscribe(buyOrderID, buyEventCh)
		defer dispatcher.Unsubscribe(buyOrderID, buyEventCh)

		placedCreditBuyOrder, err := client.PlaceOrder(ctx, creditBuyOrder)
		if err != nil {
			t.Fatalf("Failed to place credit buy order: %v", err)
		}
		assert.NotNil(t, placedCreditBuyOrder)

		//購読IDと注文IDを関連付ける
		dispatcher.RegisterOrderID(placedCreditBuyOrder.TachibanaOrderID, buyOrderID)

		// --- 買い注文の約定確認 ---
		timeout := time.After(60 * time.Second)
		for {
			select {
			case event := <-buyEventCh: // dispatcher経由でイベントを待ち受ける
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedCreditBuyOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" {
						t.Logf("Credit buy order executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						//取引区分が信用新規か確認
						assert.Equal(t, "credit_open", event.Order.TradeType, "Order trade type should be credit_open") // TradeType をチェック
						return                                                                                          // テスト成功
					} else if event.Order.Status == "4" || event.Order.Status == "5" {
						t.Fatalf("Credit buy order failed. Status: %s", event.Order.Status) //約定失敗
						return                                                              //テスト失敗
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Credit buy order execution event not received")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/credit_buy_test.go
