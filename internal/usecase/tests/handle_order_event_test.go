// internal/usecase/tests/handle_order_event_test.go
package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestTradingUsecase_HandleOrderEvent(t *testing.T) {
	testLogger := zaptest.NewLogger(t)

	// 正常系のテスト (EC イベント)
	t.Run("EC event", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, nil)
		order := &domain.Order{UUID: "order-1", Status: "filled"}
		event := &domain.OrderEvent{EventType: "EC", Order: order, Timestamp: time.Now()}

		// UpdateOrderStatus が呼ばれることを期待
		mockOrderRepo.On("UpdateOrderStatus", mock.Anything, order.UUID, order.Status).Return(nil)

		err := uc.HandleOrderEvent(context.Background(), event)
		assert.NoError(t, err)

		mockOrderRepo.AssertExpectations(t)
	})

	// 異常系のテスト (DB 更新エラー)
	t.Run("DB update error", func(t *testing.T) {
		// モッククライアントとモックリポジトリのセットアップ
		mockClient := new(tachibana.MockTachibanaClient)
		mockOrderRepo := new(usecase.MockOrderRepository)
		uc := usecase.NewTradingUsecase(mockClient, testLogger, mockOrderRepo, nil, nil)
		order := &domain.Order{UUID: "order-1", Status: "filled"}
		event := &domain.OrderEvent{EventType: "EC", Order: order}

		expectedError := errors.New("database error")
		mockOrderRepo.On("UpdateOrderStatus", mock.Anything, order.UUID, order.Status).Return(expectedError)

		err := uc.HandleOrderEvent(context.Background(), event)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)

		mockOrderRepo.AssertExpectations(t)
	})

	// 他のイベントタイプ (SS, US, NS) のテストは、ログ出力の検証が必要になるため、
	// ロギングのテストが完了した後に実装する (ここでは省略)
}
