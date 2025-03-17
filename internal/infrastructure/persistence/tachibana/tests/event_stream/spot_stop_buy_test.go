package tachibana_test

import (
	"context"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderStatus_StopBuyOrder(t *testing.T) { // テスト関数名を変更

	client := tachibana.CreateTestClient(t, nil)
	ctx := context.Background()

	t.Run("StopBuyOrderTest_GetOrderStatus", func(t *testing.T) { // テストケース名を変更
		err := client.Login(ctx, nil)
		assert.NoError(t, err)
		defer client.Logout(ctx)

		// --- 逆指値買い注文 ---
		order := &domain.Order{
			Symbol:       "7974",
			Side:         "long",
			OrderType:    "stop", // 逆指値
			Quantity:     100,
			MarketCode:   "00",
			TriggerPrice: 12000, // 適当なトリガー価格を設定（現在価格より高め）
			Price:        0,     //トリガー後、成行き
		}

		// 注文発注
		placedOrder, err := client.PlaceOrder(ctx, order)
		if err != nil {
			t.Fatalf("Failed to place order: %v", err)
		}
		assert.NotNil(t, placedOrder)
		placedOrderID := placedOrder.TachibanaOrderID
		// --- 受付確認 ---
		timeout := time.After(60 * time.Second)    // タイムアウトを設定
		orderDate := time.Now().Format("20060102") // 注文日 (当日)

		for {
			select {
			case <-timeout:
				t.Fatalf("Timeout: Order status was not '発注待ち' within the timeout period")
				return // タイムアウトしたら、テストを終了
			default: // タイムアウトしていない場合
				// GetOrderStatus を使って注文状況を確認
				statusOrder, err := client.GetOrderStatus(ctx, placedOrderID, orderDate)
				if err != nil {
					t.Logf("GetOrderStatus failed: %v, retrying...", err) // 失敗したらログ出力してリトライ
					time.Sleep(1 * time.Second)                           // 少し待ってからリトライ
					continue                                              // 次のループへ
				}

				if statusOrder.Status == "発注待ち" { // 日本語のステータスで比較
					// もし、'statusOrder.Status' が "発注待ち" であれば、
					// 注文は Tachibana API に届いているが、まだ約定待ちの状態であることを確認
					t.Logf("Order %s is pending.", placedOrderID)
					return // テスト成功
				} else if statusOrder.Status != "" {
					t.Logf("Order %s Status is : %s .", placedOrderID, statusOrder.Status)
				}
				time.Sleep(1 * time.Second) // 1秒待ってから再度確認
			}
		}
	})
}
