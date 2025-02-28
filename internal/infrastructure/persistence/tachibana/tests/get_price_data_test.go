// internal/infrastructure/persistence/tachibana/tests/get_price_data_test.go
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

func TestTachibanaClientImple_GetPriceData(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // SetupTestClient は tests パッケージで定義
	defer httpmock.DeactivateAndReset()

	t.Run("Success", func(t *testing.T) {
		// モックのレスポンスを設定 (ここでは2つの銘柄のデータを返す)

		mockResponse := `{
            "sResultCode": "0",
            "data": [
                {"sIssueCode": "7974", "sZyoukaiDay":"20240301","sHajimene":"7000", "sTakane":"7100", "sYasune":"6900", "sOwarine":"7050", "sDekidaka":"10000"},
                {"sIssueCode": "9984", "sZyoukaiDay":"20240301","sHajimene":"3000", "sTakane":"3100", "sYasune":"2900", "sOwarine":"3050", "sDekidaka":"20000"}
            ]
        }`
		httpmock.RegisterResponder("POST", "https://example.com/price", // tc.PriceURL
			httpmock.NewStringResponder(http.StatusOK, mockResponse),
		)

		issueCodes := []string{"7974", "9984"}
		priceDataList, err := client.GetPriceData(context.Background(), issueCodes)
		require.NoError(t, err)
		assert.NotNil(t, priceDataList)
		assert.Len(t, priceDataList, 2)

		// 取得したデータの検証 (例)
		assert.Equal(t, "7974", priceDataList[0].IssueCode)
		assert.Equal(t, 7000.0, priceDataList[0].Open)
		assert.Equal(t, "20240301", priceDataList[0].Date)
		assert.Equal(t, "9984", priceDataList[1].IssueCode)
		assert.Equal(t, 3000.0, priceDataList[1].Open)
		assert.Equal(t, "20240301", priceDataList[1].Date)
	})
	t.Run("API Error", func(t *testing.T) {
		// モックのレスポンスを設定 (エラーレスポンス)
		httpmock.RegisterResponder("POST", "https://example.com/price",
			httpmock.NewStringResponder(http.StatusOK, `{"sResultCode": "E999", "sResultText": "Internal Server Error"}`),
		)

		issueCodes := []string{"7974"}
		_, err := client.GetPriceData(context.Background(), issueCodes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "price data API returned an error")
	})

	t.Run("Network Error", func(t *testing.T) {
		httpmock.Reset() // モックをリセットしてネットワークエラーを発生させる

		issueCodes := []string{"7974"}
		_, err := client.GetPriceData(context.Background(), issueCodes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get price data")
	})

	t.Run("JSON Decode Error", func(t *testing.T) {
		// 不正なJSONを返すモックを設定
		httpmock.RegisterResponder("POST", "https://example.com/price",
			httpmock.NewStringResponder(http.StatusOK, `{invalid json}`),
		)

		issueCodes := []string{"7974"}
		_, err := client.GetPriceData(context.Background(), issueCodes)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get price data") // エラーメッセージを修正
		assert.Contains(t, err.Error(), "レスポンスのデコードに失敗")            //こちらのエラーメッセージを期待する
	})

	// 他の異常系のテストケース (空のリストを渡す、など) も追加
}
