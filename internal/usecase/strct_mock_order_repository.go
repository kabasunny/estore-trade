// internal/usecase/tests/mocks_test.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockOrderRepository は domain.OrderRepository のモック
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrdersBySymbolAndStatus(ctx context.Context, symbol, status string) ([]*domain.Order, error) {
	args := m.Called(ctx, symbol, status)
	return args.Get(0).([]*domain.Order), args.Error(1)
}
func (m *MockOrderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	args := m.Called(ctx, order)
	return args.Error(0)
}
func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}
