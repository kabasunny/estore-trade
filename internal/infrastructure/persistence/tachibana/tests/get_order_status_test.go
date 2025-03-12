// internal/infrastructure/persistence/tachibana/tests/get_order_status_test.go
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetOrderStatus(t *testing.T) {
	// ログイン (テストクライアントの準備)
	client := tachibana.CreateTestClient(t, nil)
	err := client.Login(context.Background(), nil)
	assert.NoError(t, err)
	defer client.Logout(context.Background())

	// 実際に存在する注文番号を指定 (テスト実行前に、デモ環境で注文を作成し、その注文番号をメモしておく)
	validOrderID := "12000006" // 例: 実際の注文番号に置き換えてください。
	orderDate := "20250312"

	t.Run("正常系: 存在する注文番号でステータスが取得できること", func(t *testing.T) {
		order, err := client.GetOrderStatus(context.Background(), validOrderID, orderDate)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, validOrderID, order.TachibanaOrderID) // 注文番号が正しいこと
		assert.NotEmpty(t, order.Status)                      // ステータスが空でないこと
		if order.Status == "executed" || order.Status == "partially_executed" {
			assert.NotZero(t, order.Quantity) //約定数量は0ではない
		}
		// 以下の確認は、domain.Order, domain.Position で long/short で統一した後
		//assert.NotEmpty(t, order.Side)        // Sideが空でないことを確認

		// 他のフィールドも必要に応じてアサーションを追加 (例: 約定単価など)
	})

	// t.Run("異常系: 存在しない注文番号でエラーになること", func(t *testing.T) {
	// 	invalidOrderID := "9999999999" // 存在しないであろう注文番号
	// 	_, err := client.GetOrderStatus(context.Background(), invalidOrderID, orderDate)
	// 	assert.Error(t, err) // エラーが発生することを期待
	// 	// エラーメッセージの検証 (必要に応じて)
	// })
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/get_order_status_test.go
