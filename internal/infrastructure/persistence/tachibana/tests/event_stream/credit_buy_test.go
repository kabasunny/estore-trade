// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_buy_test.go

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

func TestEventStreamCreditBuy(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (信用新規買い)", func(t *testing.T) {
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
		placedCreditBuyOrder, err := client.PlaceOrder(ctx, creditBuyOrder)
		if err != nil {
			t.Fatalf("Failed to place credit buy order: %v", err)
		}
		assert.NotNil(t, placedCreditBuyOrder)
		creditBuyOrderID := placedCreditBuyOrder.TachibanaOrderID

		// --- 買い注文の約定確認 ---
		timeout := time.After(60 * time.Second)
		for {
			select {
			case event := <-eventCh:
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == creditBuyOrderID {
					if event.Order.Status == "1" || event.Order.Status == "3" { // 受付済 or 一部約定
						if event.Order.FilledQuantity > 0 { //約定数量が0より大きい
							t.Logf("Credit buy order executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
							//取引区分が信用新規か確認
							assert.Equal(t, "credit_open", event.Order.TradeType, "Order trade type should be credit_open") // TradeType をチェック
							return                                                                                          // テスト成功
						}
						//assert.Equal(t, "long", event.Order.Side) //p_BBKB=3 なら long
						// 他の必要なアサーションもここに追加(約定数量、約定価格など)
					} else if event.Order.Status == "4" || event.Order.Status == "5" {
						t.Fatalf("Credit buy order failed. Status: %s", event.Order.Status) //約定失敗
						return                                                              //テスト失敗
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Credit buy order execution event not received")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestEventStreamCreditBuy
