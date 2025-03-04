package tachibana_test

import (
	"context"
	"net/http"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTachibanaClientImple_GetOrderStatus(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // テスト用のクライアントと設定を作成
	defer httpmock.DeactivateAndReset()

	t.Run("Success", func(t *testing.T) {
		mockOrderID := "12345"
		mockResponse := `{
			"sResultCode": "0",
			"sOrderNumber": "12345",
			"sOrderStatus": "1"
		}` // 必要に応じて他のフィールドも追加

		httpmock.RegisterResponder("POST", "https://example.com/request", //Login成功時に帰ってくる、RequestURLを使用
			httpmock.NewStringResponder(http.StatusOK, mockResponse),
		)

		order, err := client.GetOrderStatus(context.Background(), mockOrderID)
		require.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, mockOrderID, order.ID)
		assert.Equal(t, "1", order.Status) // "1" は仮のステータス値
		// 他のフィールドも必要に応じて検証
	})

	t.Run("API Error", func(t *testing.T) {
		mockOrderID := "12345"
		mockResponse := `{"sResultCode": "E999", "sResultText": "Order not found"}` // エラーレスポンス

		httpmock.RegisterResponder("POST", "https://example.com/request", //Login成功時に帰ってくる、RequestURLを使用
			httpmock.NewStringResponder(http.StatusOK, mockResponse), // 200 OK でエラーレスポンスを返す
		)

		_, err := client.GetOrderStatus(context.Background(), mockOrderID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "order status API returned an error")
	})

	t.Run("Network Error", func(t *testing.T) {
		mockOrderID := "12345"

		httpmock.Reset() // Responder をクリアしてネットワークエラーを発生させる

		_, err := client.GetOrderStatus(context.Background(), mockOrderID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "get order status failed")
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		mockOrderID := "12345"
		mockResponse := `{invalid json}` // 不正なJSON

		httpmock.RegisterResponder("POST", "https://example.com/request",
			httpmock.NewStringResponder(http.StatusOK, mockResponse),
		)

		_, err := client.GetOrderStatus(context.Background(), mockOrderID)
		require.Error(t, err)                            // エラーが発生することを期待
		assert.Contains(t, err.Error(), "レスポンスのデコードに失敗") // エラーメッセージを修正
	})
}
