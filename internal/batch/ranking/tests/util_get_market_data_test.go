// util_get_market_data_test.go
package ranking_test

import (
	"context"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMarketData(t *testing.T) {
	// モッククライアントの作成
	mockClient := new(tachibana.MockTachibanaClient)

	// GetRequestURL メソッドのモック設定
	mockRequestURL := "mocked_request_url"
	mockClient.On("GetRequestURL").Return(mockRequestURL, nil) // 戻り値を設定

	// テスト対象の関数を呼び出す
	issueCodes := []string{"1301", "1305"}
	marketData, err := ranking.GetMarketData(context.Background(), mockClient, issueCodes)

	// アサーション
	assert.NoError(t, err)
	assert.Empty(t, marketData) // marketData が空であることを確認

	// モックが期待通りに呼び出されたかを確認
	mockClient.AssertExpectations(t)
}
