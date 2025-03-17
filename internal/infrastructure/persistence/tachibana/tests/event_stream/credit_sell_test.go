// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_test.go

package tachibana_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher のインポート
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestEventStreamCreditSell(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("イベントストリームのテスト (信用新規売り)", func(t *testing.T) {
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

		// --- 信用新規売り注文 ---
		sellOrder := &domain.Order{
			Symbol:     "7974",   // 例: 任天堂
			Side:       "short",  // 売り
			OrderType:  "market", // 成行 (または "limit")
			Quantity:   100,
			MarketCode: "00",          // 東証
			TradeType:  "credit_open", // 信用新規
			//ExecutionType: "market",      // 今回は執行条件をmarketで固定  <-- 不要なので削除
			// Price:      13000,      // 指値の場合 (例)
		}

		// 売り注文用のチャネルと購読ID
		sellEventCh := make(chan *domain.OrderEvent, 10)
		sellOrderID := "sellOrder"
		dispatcher.Subscribe(sellOrderID, sellEventCh)
		defer dispatcher.Unsubscribe(sellOrderID, sellEventCh)

		placedSellOrder, err := client.PlaceOrder(ctx, sellOrder)
		if err != nil {
			t.Fatalf("Failed to place credit sell order: %v", err)
		}
		assert.NotNil(t, placedSellOrder)
		// 購読IDと注文IDを関連付ける
		dispatcher.RegisterOrderID(placedSellOrder.TachibanaOrderID, sellOrderID)

		// --- 売り注文の約定確認 ---
		timeout := time.After(60 * time.Second) // タイムアウトを設定 (60秒)
		for {
			select {
			case event := <-sellEventCh: // dispatcher 経由でイベント受信
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedSellOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" { //全部約定
						t.Logf("Credit sell order executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						assert.Equal(t, "credit_open", event.Order.TradeType, "Order trade type should be credit_open") // TradeType を確認
						return                                                                                          // テスト成功
					} else if event.Order.Status == "4" || event.Order.Status == "5" {
						t.Fatalf("Credit sell order failed. Status: %s", event.Order.Status) //約定失敗
						return
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Credit sell order execution event not received")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_test.go

// 信用新規売り注文 (成行) を発注。
// イベントストリーム経由で約定イベント (EC) を受信。
// 受信したイベントから、注文番号、ステータス、約定数量、取引区分などを確認。
// 注文が正常に約定していれば、テスト成功
