// internal/infrastructure/persistence/tachibana/tests/tachibana_client_mock_test.go
package tachibana_test // tests -> tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana" // tachibana パッケージをインポート

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockTachibanaClient(t *testing.T) {
	// モッククライアントのインスタンスを作成
	mockClient := new(tachibana.MockTachibanaClient)

	// Login メソッドのモックを設定 (引数なし、エラーなし)
	mockClient.On("Login", mock.Anything, mock.Anything).Return(nil)

	// テスト対象のメソッドを呼び出す
	err := mockClient.Login(context.Background(), nil)

	// 期待される結果を検証
	assert.NoError(t, err) // エラーが発生しないことを確認

	// モックが期待通りに呼び出されたか検証
	mockClient.AssertExpectations(t)
}
