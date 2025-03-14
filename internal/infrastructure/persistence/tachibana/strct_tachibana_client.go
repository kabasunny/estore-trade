// internal/infrastructure/persistence/tachibana/strct_tachibana_client.go
package tachibana

import (
	"estore-trade/internal/domain"
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TachibanaClientImple 構造体の定義
type TachibanaClientImple struct {
	baseURL  *url.URL
	userID   string // 立花証券ユーザーID
	password string // 立花証券パスワード
	// apiKey     string
	secret     string
	logger     *zap.Logger
	loggined   bool
	requestURL string       // キャッシュする仮想URL（REQUEST)
	masterURL  string       // キャッシュする仮想URL（Master)
	priceURL   string       // キャッシュする仮想URL（Price)
	eventURL   string       // キャッシュする仮想URL（EVENT)
	expiry     time.Time    // 仮想URLの有効期限
	mu         sync.RWMutex // 排他制御用
	pNo        int64        // p_no の連番管理用
	pNoMu      sync.Mutex   // pNo の排他制御用

	targetIssueCodes   []string
	targetIssueCodesMu sync.RWMutex // 多分いらんので最後に消す
	masterData         *domain.MasterData

	//HTTPClient HTTPClient // HTTPClientインターフェースのフィールドを削除
}
