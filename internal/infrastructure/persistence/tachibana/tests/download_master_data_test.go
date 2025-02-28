// internal/infrastructure/persistence/tachibana/tests/download_master_data_test.go
package tachibana_test

import (
	"context"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTachibanaClientImple_DownloadMasterData(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // SetupTestClient は tests パッケージで定義
	defer httpmock.DeactivateAndReset()

	t.Run("Success", func(t *testing.T) {
		// モックのレスポンスを設定 (ここでは簡略化のため、一部のマスタデータのみ)
		mockResponse := `{
            "CLMSystemStatus": {
                "sCLMID": "CLMSystemStatus",
                "sSystemStatusKey": "001",
                "sLoginKyokaKubun": "1",
                "sSystemStatus": "1"
            },
            "CLMDateZyouhou": {
                "sCLMID": "CLMDateZyouhou",
                "sDayKey": "001",
                "sTheDay": "20231101"
            },
            "CLMEventDownloadComplete": {}
        }`
		httpmock.RegisterResponder("POST", "https://example.com/master", //tc.MasterURL
			httpmock.NewStringResponder(http.StatusOK, mockResponse),
		)

		masterData, err := client.DownloadMasterData(context.Background())
		require.NoError(t, err)
		assert.NotNil(t, masterData)

		// 取得したマスタデータの検証 (一部)
		assert.Equal(t, "1", masterData.SystemStatus.SystemState)
		assert.Equal(t, "20231101", masterData.DateInfo.TheDay)

		// 他のマスタデータについても、必要に応じて検証を追加
	})
	t.Run("API Error", func(t *testing.T) {
		// エラーレスポンスを返すモックを設定
		httpmock.RegisterResponder("POST", "https://example.com/master", //
			httpmock.NewStringResponder(http.StatusInternalServerError, ""),
		)

		_, err := client.DownloadMasterData(context.Background())
		require.Error(t, err) // エラーが発生することを期待
	})
	t.Run("JSON Decode Error", func(t *testing.T) {
		// 不正なJSONを返すモックを設定
		httpmock.RegisterResponder("POST", "https://example.com/master", //tc.MasterURL
			httpmock.NewStringResponder(http.StatusOK, `{invalid json}`),
		)

		_, err := client.DownloadMasterData(context.Background())
		require.Error(t, err)                            // エラーが発生することを期待
		assert.Contains(t, err.Error(), "レスポンスのデコードに失敗") // エラーメッセージを修正
	})

	// 他の異常系のテストケース (ネットワークエラーなど) も追加
}
