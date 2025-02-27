// internal/usecase/tests/cancel_order_test.go
package usecase_test // usecase -> usecase_test

import (
	"context"
	"errors"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase" // usecase パッケージをインポート

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository, MockAccountRepository はこのファイル内にコピー (または place_order_test.go から import)

func TestTradingUsecase_CancelOrder(t *testing.T) {

	orderID := "test-order-id"

	// 正常系のテスト
	t.Run("valid order id", func(t *testing.T) {
		//モックのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		uc := usecase.NewTradingUsecase(mockClient, nil, nil, nil, nil) // usecase. を追加
		mockClient.On("CancelOrder", mock.Anything, orderID).Return(nil)

		err := uc.CancelOrder(context.Background(), orderID)
		assert.NoError(t, err)

		mockClient.AssertExpectations(t)
	})

	// 異常系のテスト (APIエラー)
	t.Run("API error", func(t *testing.T) {
		//モックのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		uc := usecase.NewTradingUsecase(mockClient, nil, nil, nil, nil)
		expectedError := errors.New("tachibana API error")
		mockClient.On("CancelOrder", mock.Anything, orderID).Return(expectedError)

		err := uc.CancelOrder(context.Background(), orderID)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})
}
