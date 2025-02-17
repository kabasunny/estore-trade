// internal/infrastructure/persistence/tachibana/tachibana_client.go
package tachibana

import (
	"context"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
)

// TachibanaClient インターフェースは、立花証券のAPIとやり取りするためのメソッドを定義
type TachibanaClient interface {
	// APIに対してログインし、ユーザーIDとパスワードを使用して必要な認証情報を取得し、成功した場合、APIとやり取りするためのリクエストURLを返す
	Login(ctx context.Context, cfg *config.Config) (bool, error)

	// 新しい株式注文を立花証券のAPIに対して行い、 注文が成功した場合、注文情報を含む domain.Order オブジェクトを返す
	PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error)

	// 指定された注文IDに基づいて注文のステータスを取得し、注文のステータス情報を含む domain.Order オブジェクトを返す
	GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error)

	// 指定された注文IDに基づいて注文のキャンセルを行い、キャンセルが成功した場合はエラーを返さない
	CancelOrder(ctx context.Context, orderID string) error

	ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error) // OrderEventチャネル

	DownloadMasterData(ctx context.Context) error // マスタデータダウンロード

	GetSystemStatus() SystemStatus // SystemStatus を返す
	GetDateInfo() DateInfo
	GetCallPrice(unitNumber string) (CallPrice, bool)
	GetIssueMaster(issueCode string) (IssueMaster, bool)

	// 各URLを取得するためのメソッド
	GetRequestURL() (string, error)
	GetMasterURL() (string, error)
	GetPriceURL() (string, error)
	GetEventURL() (string, error)
}
