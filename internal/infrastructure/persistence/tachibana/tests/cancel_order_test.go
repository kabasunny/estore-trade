// internal/infrastructure/persistence/tachibana/tests/cancel_order_test.go
package tachibana_test

import (
	"context"
	"net/http"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana" // tachibana パッケージをインポート

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTachibanaClientImple_CancelOrder(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // tachibana.SetupTestClient(t) を呼び出す
	defer httpmock.DeactivateAndReset()

	t.Run("Success", func(t *testing.T) {
		mockOrderID := "12345"

		// モックのレスポンスを設定
		httpmock.RegisterResponder("POST", "https://example.com/request/cancel", //RequestURLに/cancelを追加
			httpmock.NewStringResponder(http.StatusOK, `{"sResultCode": "0"}`),
		)

		err := client.CancelOrder(context.Background(), mockOrderID)
		require.NoError(t, err)
	})

	t.Run("API Error", func(t *testing.T) {
		mockOrderID := "12345"

		// モックのレスポンスを設定 (エラーレスポンス)
		httpmock.RegisterResponder("POST", "https://example.com/request/cancel", //RequestURLに/cancelを追加
			httpmock.NewStringResponder(http.StatusOK, `{"sResultCode": "E999", "sResultText": "Order not found"}`),
		)

		err := client.CancelOrder(context.Background(), mockOrderID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cancel order API returned an error")
	})

	t.Run("Network Error", func(t *testing.T) {
		mockOrderID := "12345"

		httpmock.Reset() // モックをリセットしてネットワークエラーを発生させる

		err := client.CancelOrder(context.Background(), mockOrderID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cancel order failed")
	})
}
