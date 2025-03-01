// internal/infrastructure/persistence/tachibana/iface_client.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain"
)

type TachibanaClient interface {
	Login(ctx context.Context, cfg interface{}) error
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)
	CancelOrder(ctx context.Context, orderID string) error
	ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error)
	GetRequestURL() (string, error)
	// GetMasterURL() (string, error)
	// GetPriceURL() (string, error)
	GetEventURL() (string, error)
	DownloadMasterData(ctx context.Context) (*domain.MasterData, error)
	GetSystemStatus() domain.SystemStatus // システムステータス取得
	// GetDateInfo() domain.DateInfo         // 日付情報取得
	// GetCallPrice(unitNumber string) (domain.CallPrice, bool)
	GetIssueMaster(issueCode string) (domain.IssueMaster, bool)
	// GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool)
	// GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool)
	// GetOperationStatusKabu(listedMarket string, unit string) (domain.OperationStatusKabu, bool)
	CheckPriceIsValid(issueCode string, price float64, isNextDay bool) (bool, error)
	// SetTargetIssues(ctx context.Context, issueCodes []string) error
	// GetPriceData(ctx context.Context, issueCodes []string) ([]domain.PriceData, error)
	// GetMasterData() *domain.MasterData
}
