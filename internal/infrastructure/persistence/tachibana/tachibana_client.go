// internal/infrastructure/persistence/tachibana/tachibana.go
package tachibana

import (
	"context"
	"estore-trade/internal/domain" // OrderEvent構造体を使用するため
)

// TachibanaClient インターフェース (メソッドのシグネチャを定義)
type TachibanaClient interface {
	Login(ctx context.Context, cfg interface{}) error                                    // ログイン
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)          // 注文
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)           // 注文状態取得
	CancelOrder(ctx context.Context, orderID string) error                               // 注文取消
	ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error)           // イベントストリーム接続
	GetRequestURL() (string, error)                                                      // リクエストURL取得
	GetMasterURL() (string, error)                                                       // マスタURL取得
	GetPriceURL() (string, error)                                                        // プライスURL取得
	GetEventURL() (string, error)                                                        // イベントURL取得
	DownloadMasterData(ctx context.Context) error                                        // マスタデータダウンロード
	GetSystemStatus() SystemStatus                                                       // システムステータス取得
	GetDateInfo() DateInfo                                                               // 日付情報取得
	GetCallPrice(unitNumber string) (CallPrice, bool)                                    // 呼値取得
	GetIssueMaster(issueCode string) (IssueMaster, bool)                                 // 株式銘柄マスタ取得
	GetIssueMarketMaster(issueCode, marketCode string) (IssueMarketMaster, bool)         // 株式銘柄市場マスタ取得
	GetIssueMarketRegulation(issueCode, marketCode string) (IssueMarketRegulation, bool) // 株式銘柄市場別規制取得
	GetOperationStatusKabu(listedMarket string, unit string) (OperationStatusKabu, bool) // 運用ステータス（株）取得
	CheckPriceIsValid(issueCode string, price float64, isNextDay bool) (bool, error)     // 追加: 呼値チェック
	SetTargetIssues(ctx context.Context, issueCodes []string) error                      // 追加: ターゲット銘柄を設定
}
