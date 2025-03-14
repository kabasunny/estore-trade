// internal/usecase/tests/get_order_status_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTradingUsecase_GetOrderStatus(t *testing.T) {

	orderID := "test-order-id"

	// 正常系のテスト
	t.Run("valid order id", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		uc := usecase.NewTradingUsecase(mockClient, nil, nil, nil, nil)
		expectedOrder := &domain.Order{UUID: orderID, Status: "filled"}
		mockClient.On("GetOrderStatus", mock.Anything, orderID).Return(expectedOrder, nil)

		order, err := uc.GetOrderStatus(context.Background(), orderID)
		assert.NoError(t, err)
		assert.NotNil(t, order)
		assert.Equal(t, expectedOrder, order)

		mockClient.AssertExpectations(t)
	})

	// 異常系のテスト (APIエラー)
	t.Run("API error", func(t *testing.T) {
		mockClient := new(tachibana.MockTachibanaClient)
		uc := usecase.NewTradingUsecase(mockClient, nil, nil, nil, nil)
		expectedError := errors.New("tachibana API error")
		mockClient.On("GetOrderStatus", mock.Anything, orderID).Return((*domain.Order)(nil), expectedError)

		order, err := uc.GetOrderStatus(context.Background(), orderID)
		assert.Error(t, err)
		assert.Nil(t, order)
		assert.Equal(t, expectedError, err)

		mockClient.AssertExpectations(t)
	})
}
