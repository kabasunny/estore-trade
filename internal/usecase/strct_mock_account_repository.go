// internal/usecase/tests/mocks_test.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockAccountRepository は domain.AccountRepository のモック
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) GetAccountByUserID(ctx context.Context, userID string) (*domain.Account, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Account), args.Error(1)
}

func (m *MockAccountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}
func (m *MockAccountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}
