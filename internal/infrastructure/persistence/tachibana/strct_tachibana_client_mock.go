// internal/infrastructure/persistence/tachibana/strct_tachibana_client_mock.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockTachibanaClient は TachibanaClient インターフェースのモックです。
type MockTachibanaClient struct {
	mock.Mock
}

func (m *MockTachibanaClient) Login(ctx context.Context, cfg interface{}) error {
	args := m.Called(ctx, cfg)
	return args.Error(0)
}

func (m *MockTachibanaClient) Logout(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockTachibanaClient) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	args := m.Called(ctx, order)
	if err := args.Error(1); err != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockTachibanaClient) GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error) {
	args := m.Called(ctx, orderID, orderDate) // 引数を修正
	return args.Get(0).(*domain.Order), args.Error(1)
}

func (m *MockTachibanaClient) CancelOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockTachibanaClient) ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error) {
	args := m.Called(ctx)
	return args.Get(0).(<-chan *domain.OrderEvent), args.Error(1)
}

func (m *MockTachibanaClient) DownloadMasterData(ctx context.Context) (*domain.MasterData, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.MasterData), args.Error(1)
}

func (m *MockTachibanaClient) GetSystemStatus(ctx context.Context) domain.SystemStatus {
	args := m.Called(ctx) // 引数を追加
	return args.Get(0).(domain.SystemStatus)
}

func (m *MockTachibanaClient) GetIssueMaster(ctx context.Context, issueCode string) (domain.IssueMaster, bool) {
	args := m.Called(ctx, issueCode) // 引数を追加
	return args.Get(0).(domain.IssueMaster), args.Bool(1)
}

func (m *MockTachibanaClient) CheckPriceIsValid(ctx context.Context, issueCode string, price float64, isNextDay bool) (bool, error) {
	args := m.Called(ctx, issueCode, price, isNextDay) // 引数を修正
	return args.Bool(0), args.Error(1)
}

func (m *MockTachibanaClient) GetPositions(ctx context.Context) ([]domain.Position, error) {
	args := m.Called(ctx) // 引数を追加
	return args.Get(0).([]domain.Position), args.Error(1)
}

func (m *MockTachibanaClient) SetTargetIssues(ctx context.Context, issueCodes []string) error {
	args := m.Called(ctx, issueCodes)
	return args.Error(0)
}

// 以下は、tachibana.TachibanaClient インターフェースに存在しないため、削除
/*
func (m *MockTachibanaClient) GetRequestURL() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTachibanaClient) GetMasterURL() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockTachibanaClient) GetPriceURL() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}
func (m *MockTachibanaClient) GetEventURL() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockTachibanaClient) GetDateInfo() domain.DateInfo {
	args := m.Called()
	return args.Get(0).(domain.DateInfo) //
}
func (m *MockTachibanaClient) GetCallPrice(unitNumber string) (domain.CallPrice, bool) {
	args := m.Called(unitNumber)
	return args.Get(0).(domain.CallPrice), args.Bool(1) //
}

func (m *MockTachibanaClient) GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool) {
	args := m.Called(issueCode, marketCode)
	return args.Get(0).(domain.IssueMarketMaster), args.Bool(1) //
}
func (m *MockTachibanaClient) GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool) {
	args := m.Called(issueCode, marketCode)
	return args.Get(0).(domain.IssueMarketRegulation), args.Bool(1) //
}
func (m *MockTachibanaClient) GetOperationStatusKabu(listedMarket string, unit string) (domain.OperationStatusKabu, bool) {
	args := m.Called(listedMarket, unit)
	return args.Get(0).(domain.OperationStatusKabu), args.Bool(1) //
}

func (m *MockTachibanaClient) GetPriceData(ctx context.Context, issueCodes []string) ([]domain.PriceData, error) {
	args := m.Called(ctx, issueCodes)
	return args.Get(0).([]domain.PriceData), args.Error(1)
}

func (m *MockTachibanaClient) GetMasterData() *domain.MasterData {
	args := m.Called()
	return args.Get(0).(*domain.MasterData)
}
*/
