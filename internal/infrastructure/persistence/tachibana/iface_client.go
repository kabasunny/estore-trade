// internal/infrastructure/persistence/tachibana/tachibana.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain" // OrderEvent構造体を使用するため
)

// TachibanaClient インターフェース (メソッドのシグネチャを定義)
// 基本的には、値渡しを優先し、明確な理由がある場合にのみポインタ渡しを選択するという方針が良い
// 今回のケースでは、SystemStatus と IssueMaster は値渡し、Order はポインタ渡し、MasterData は (現状ではどちらでも良いが) ポインタ渡しとしておくのが、現時点での最適解と考えられる
type TachibanaClient interface {
	Login(ctx context.Context, cfg interface{}) error                                           // ログイン
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)                 // 注文
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)                  // 注文状態取得
	CancelOrder(ctx context.Context, orderID string) error                                      // 注文取消
	ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error)                  // イベントストリーム接続
	GetRequestURL() (string, error)                                                             // リクエストURL取得
	GetMasterURL() (string, error)                                                              // マスタURL取得
	GetPriceURL() (string, error)                                                               // プライスURL取得
	GetEventURL() (string, error)                                                               // イベントURL取得
	DownloadMasterData(ctx context.Context) (*domain.MasterData, error)                         // マスタデータダウンロード
	GetSystemStatus() domain.SystemStatus                                                       // システムステータス取得
	GetDateInfo() domain.DateInfo                                                               // 日付情報取得
	GetCallPrice(unitNumber string) (domain.CallPrice, bool)                                    // 呼値取得
	GetIssueMaster(issueCode string) (domain.IssueMaster, bool)                                 // 株式銘柄マスタ取得
	GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool)         // 株式銘柄市場マスタ取得
	GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool) // 株式銘柄市場別規制取得
	GetOperationStatusKabu(listedMarket string, unit string) (domain.OperationStatusKabu, bool) // 運用ステータス（株）取得
	CheckPriceIsValid(issueCode string, price float64, isNextDay bool) (bool, error)            // 呼値チェック
	SetTargetIssues(ctx context.Context, issueCodes []string) error                             // ターゲット銘柄を設定
	GetPriceData(ctx context.Context, issueCodes []string) ([]domain.PriceData, error)          // 銘柄リストの株価データを取得
	//GetAllIssueCodes(ctx context.Context) ([]string, error)                              // 全銘柄リストを設定
	GetMasterData() *domain.MasterData // masterDataManagerのゲッター
}
