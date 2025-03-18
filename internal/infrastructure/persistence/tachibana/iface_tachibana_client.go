// internal/infrastructure/persistence/tachibana/iface_tachibana_client.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain"
)

// APIからデータを取得するためのクライアント
type TachibanaClient interface {
	Login(ctx context.Context, cfg interface{}) error
	Logout(ctx context.Context) error

	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)
	GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error)
	CancelOrder(ctx context.Context, orderID string) error
	ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error)
	// GetRequestURL() (string, error)
	// GetMasterURL() (string, error)
	// GetPriceURL() (string, error)
	// GetEventURL() (string, error)
	DownloadMasterData(ctx context.Context) (*domain.MasterData, error)
	GetSystemStatus(ctx context.Context) domain.SystemStatus // システムステータス取得
	// GetDateInfo() domain.DateInfo         // 日付情報取得
	// GetCallPrice(unitNumber string) (domain.CallPrice, bool)
	GetIssueMaster(ctx context.Context, issueCode string) (domain.IssueMaster, bool)
	// GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool)
	// GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool)
	// GetOperationStatusKabu(listedMarket string, unit string) (domain.OperationStatusKabu, bool)
	CheckPriceIsValid(ctx context.Context, issueCode string, price float64, isNextDay bool) (bool, error)
	SetTargetIssues(ctx context.Context, issueCodes []string) error
	// GetPriceData(ctx context.Context, issueCodes []string) ([]domain.PriceData, error)
	// GetMasterData() *domain.MasterData

	GetPositions(ctx context.Context) ([]domain.Position, error)
}
