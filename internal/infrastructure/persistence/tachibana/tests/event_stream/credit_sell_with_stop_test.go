// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_with_stop_test.go
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

func TestCreditSellWithStopCombined_GetOrderStatus(t *testing.T) { // 関数名を変更
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("信用新規成行売り + 逆指値買い - 同時注文 (EventStream + GetOrderStatus)", func(t *testing.T) { // テストケース名を変更
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

		// --- 信用新規成行売り + 逆指値買い注文 ---  <-- コメントは残す
		sellOrder := &domain.Order{
			Symbol:       "7974", // 例: 任天堂
			Side:         "short",
			OrderType:    "stop",        // 通常 + 逆指値
			TradeType:    "credit_open", // 信用新規
			Quantity:     100,
			MarketCode:   "00",    // 東証
			Price:        0,       // 成行
			TriggerPrice: 11000.0, // 例: 現在価格が10000円として、11000円以上になったら
		}

		// 売り注文用のチャネルと購読ID (イベント受信用)
		sellEventCh := make(chan *domain.OrderEvent, 10)
		sellOrderID := "sellOrder" // 購読ID。注文IDとは異なるユニークなものにする
		dispatcher.Subscribe(sellOrderID, sellEventCh)
		defer dispatcher.Unsubscribe(sellOrderID, sellEventCh)

		placedSellOrder, err := client.PlaceOrder(ctx, sellOrder)
		if err != nil {
			t.Fatalf("Failed to place sell order: %v", err)
		}
		assert.NotNil(t, placedSellOrder)
		// 購読IDと注文IDを関連づける
		dispatcher.RegisterOrderID(placedSellOrder.TachibanaOrderID, sellOrderID)

		// --- 注文状態の確認（逆指値なので、最初は受付済） ---  <-- コメント修正
		timeout := time.After(60 * time.Second)
	OrderLoop:
		for {
			select {
			case event := <-sellEventCh: // dispatcher経由でイベント受信
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				//イベントでの約定確認
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedSellOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" { //全部約定
						t.Logf("Sell order fully executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						// テスト成功 (売り注文の約定を確認)
						break OrderLoop //約定したので、ループを抜ける
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Sell order failed. Status: %s", event.Order.Status)
						return
					}
					// Statusが0や1の場合は何もせず、次のイベントを待つ（ループを継続）
					// ログで状況確認(status=1(一部約定)の時もログ出力)
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)

				}
			case <-timeout:
				t.Fatalf("Timeout: Sell order execution event not received. Last status: %s", "Unknown") //timeout時のstatus不明
				return
			default:
				//ここでは、GetOrderStatus は呼ばない
				time.Sleep(1 * time.Second)
			}
		}
	})
}

//  go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_with_stop_test.go
