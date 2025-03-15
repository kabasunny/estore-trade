// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_buy_sell_test.go

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

func TestCreditBuyAndRepayWithPositionID(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("信用新規買い -> ポジションIDで特定して返済売り", func(t *testing.T) {
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
		buyOrder := &domain.Order{
			Symbol:        "7974", // 例: 任天堂
			Side:          "long",
			OrderType:     "market",
			Quantity:      100,
			MarketCode:    "00",          // 東証
			TradeType:     "credit_open", // 信用新規
			ExecutionType: "market",      // 今回は執行条件をmarketで固定
			// Price:      13000,      // 指値の場合 (例)
		}
		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place credit buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)
		buyOrderID := placedBuyOrder.TachibanaOrderID

		// --- 買い注文の約定確認 ---
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
						t.Logf("Credit buy order executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						buyEvent = event
						break BuyOrderLoop // 買い注文の約定確認ループを抜ける
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Credit buy order execution event not received")
				return
			}
		}

		// --- ポジション情報取得  ---
		//time.Sleep(5 * time.Second) // GetPositions の結果が反映されるまで待機
		positions, err := client.GetPositions(ctx) //ポジションリストを取得
		if err != nil {
			t.Fatalf("Failed to get positions: %v", err)
		}

		// --- 返済注文用ポジション特定 ---
		var positionToRepay *domain.Position
		buyEventDate := buyEvent.Order.CreatedAt.Format("20060102") //追加
		for i, p := range positions {
			// 銘柄、Side、建日が一致するポジションを探す
			if p.Symbol == buyEvent.Order.Symbol && p.Side == "long" && p.OpenDate.Format("20060102") == buyEventDate && p.Quantity >= buyEvent.Order.FilledQuantity {
				positionToRepay = &positions[i] // ポインタを代入
				break
			}
		}

		if positionToRepay == nil {
			t.Fatal("Failed to find matching position for repayment")
			return
		}
		t.Logf("positionToRepay: %+v\n", *positionToRepay) // 確認

		// --- 信用返済売り注文 (ポジションIDを指定) ---
		sellOrder := &domain.Order{
			Symbol:        positionToRepay.Symbol,   // 買い注文と同じ銘柄
			Side:          "short",                  // 売り
			OrderType:     "credit_close_market",    // 信用返済 (成行)
			Quantity:      positionToRepay.Quantity, // 買い注文の約定数量
			MarketCode:    "00",                     // 買い注文と同じ市場
			TradeType:     "credit_close",           // 追加: 信用返済
			ExecutionType: "market",                 // 追加: 成行
			Positions: []domain.Position{
				{
					ID:       positionToRepay.ID, // ここでポジションIDを使用
					Quantity: positionToRepay.Quantity,
				},
			},
		}

		placedSellOrder, err := client.PlaceOrder(ctx, sellOrder)
		if err != nil {
			t.Fatalf("Failed to place sell order: %v", err)
		}
		assert.NotNil(t, placedSellOrder)
		sellOrderID := placedSellOrder.TachibanaOrderID

		// --- 売り注文の約定確認 ---
		timeout = time.After(60 * time.Second) // タイムアウトをリセット
		for {
			select {
			case event := <-eventCh:
				if event == nil { //nilチェック
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == sellOrderID {
					if event.Order.Status == "1" || event.Order.Status == "3" {
						if event.Order.FilledQuantity > 0 {
							t.Logf("Sell order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
							return // テスト成功 (売り注文の約定を確認)
						}
					}
				}
			case <-timeout:
				t.Fatalf("Timeout: Sell order execution event not received")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream -run TestCreditBuyAndRepayWithPositionID

// 信用新規買い注文を発注
// 約定イベントを待機 (買い)
// GetPositions で建玉情報を取得
// 建玉情報から、返済対象のポジションを特定
// 特定したポジションIDを使って、信用返済売り注文を発注
// 約定イベントを待機 (売り)
// 約定確認後、テスト終了 (ログアウト)
