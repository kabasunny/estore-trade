// internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_buy_test.go

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

func TestCreditSellAndRepayWithPositionID(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("信用新規売り -> ポジションIDで特定して返済買い", func(t *testing.T) {
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

		// --- 信用新規売り注文 ---
		sellOrder := &domain.Order{
			Symbol:        "7974", // 例: 任天堂
			Side:          "short",
			OrderType:     "market",
			Quantity:      100,
			MarketCode:    "00",          // 東証
			TradeType:     "credit_open", // 信用新規
			ExecutionType: "market",      // 今回は執行条件をmarketで固定
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

		//購読IDと注文IDを関連付ける
		dispatcher.RegisterOrderID(placedSellOrder.TachibanaOrderID, sellOrderID)

		// --- 売り注文の約定確認 ---
		var sellEvent *domain.OrderEvent
		timeout := time.After(60 * time.Second)
	SellOrderLoop:
		for {
			select {
			case event := <-sellEventCh: // dispatcher 経由でイベント受信
				if event == nil {
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedSellOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" {
						t.Logf("Credit sell order executed. Status: %s, Executed Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						sellEvent = event
						break SellOrderLoop // 売り注文の約定確認ループを抜ける
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Credit sell order failed. Status: %s", event.Order.Status)
						return
					}
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Credit sell order execution event not received")
				return
			}
		}

		// --- ポジション情報取得  ---
		positions, err := client.GetPositions(ctx) //ポジションリストを取得
		if err != nil {
			t.Fatalf("Failed to get positions: %v", err)
		}

		// --- 返済注文用ポジション特定 ---
		var positionToRepay *domain.Position
		sellEventDate := sellEvent.Order.CreatedAt.Format("20060102") //追加
		for i, p := range positions {
			// 銘柄、Side、建日が一致するポジションを探す
			if p.Symbol == sellEvent.Order.Symbol && p.Side == "short" && p.OpenDate.Format("20060102") == sellEventDate && p.Quantity >= sellEvent.Order.FilledQuantity {
				positionToRepay = &positions[i] // ポインタを代入
				break
			}
		}

		if positionToRepay == nil {
			t.Fatal("Failed to find matching position for repayment")
			return
		}
		t.Logf("positionToRepay: %+v\n", *positionToRepay) // 確認

		// --- 信用返済買い注文 (ポジションIDを指定) ---
		buyOrder := &domain.Order{
			Symbol:        positionToRepay.Symbol,   // 売り注文と同じ銘柄
			Side:          "long",                   // 買い
			OrderType:     "credit_close_market",    // 信用返済 (成行)
			Quantity:      positionToRepay.Quantity, // 売り注文の約定数量
			MarketCode:    "00",                     // 売り注文と同じ市場
			TradeType:     "credit_close",           // 追加: 信用返済
			ExecutionType: "market",                 // 追加: 成行
			Positions: []domain.Position{
				{
					ID:       positionToRepay.ID, // ここでポジションIDを使用
					Quantity: positionToRepay.Quantity,
				},
			},
		}
		//買い注文用のチャネルと購読ID
		buyEventCh := make(chan *domain.OrderEvent, 10) // バッファを設定
		buyOrderID := "buyOrder"                        //買い注文用の購読ID
		dispatcher.Subscribe(buyOrderID, buyEventCh)
		defer dispatcher.Unsubscribe(buyOrderID, buyEventCh)

		placedBuyOrder, err := client.PlaceOrder(ctx, buyOrder)
		if err != nil {
			t.Fatalf("Failed to place buy order: %v", err)
		}
		assert.NotNil(t, placedBuyOrder)

		//購読IDと注文IDを関連付ける
		dispatcher.RegisterOrderID(placedBuyOrder.TachibanaOrderID, buyOrderID)

		// --- 買い注文の約定確認 ---
		timeout = time.After(60 * time.Second) // タイムアウトをリセット
		for {
			select {
			case event := <-buyEventCh: //dispatcher経由でイベント受信
				if event == nil { //nilチェック
					continue
				}
				fmt.Printf("Received event: %+v\n", event)
				if event.EventType == "EC" && event.Order != nil && event.Order.TachibanaOrderID == placedBuyOrder.TachibanaOrderID {
					if event.Order.ExecutionStatus == "2" {
						t.Logf("Buy order executed. Status: %s, Quantity: %d", event.Order.Status, event.Order.FilledQuantity)
						return // テスト成功 (買い注文の約定を確認)
					} else if event.Order.Status == "4" || event.Order.Status == "5" { //注文失敗
						t.Fatalf("Credit buy order failed. Status: %s", event.Order.Status)
						return
					}
					t.Logf("Waiting for full execution. Current Status: %s, FilledQuantity: %d", event.Order.Status, event.Order.FilledQuantity)
				}
			case <-timeout:
				t.Fatalf("Timeout: Buy order execution event not received")
				return
			}
		}
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/event_stream/credit_sell_buy_test.go
