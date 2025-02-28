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
	BaseURL    *url.URL
	ApiKey     string
	Secret     string
	Logger     *zap.Logger
	Loggined   bool
	RequestURL string       // キャッシュする仮想URL（REQUEST)
	MasterURL  string       // キャッシュする仮想URL（Master)
	PriceURL   string       // キャッシュする仮想URL（Price)
	EventURL   string       // キャッシュする仮想URL（EVENT)
	Expiry     time.Time    // 仮想URLの有効期限
	Mu         sync.RWMutex // 排他制御用
	PNo        int64        // p_no の連番管理用
	PNoMu      sync.Mutex   // pNo の排他制御用

	// マスタデータ保持用 (必要最低限に絞り込み)
	systemStatus             domain.SystemStatus
	dateInfo                 domain.DateInfo
	callPriceMap             map[string]domain.CallPrice
	issueMap                 map[string]domain.IssueMaster
	issueMarketMap           map[string]map[string]domain.IssueMarketMaster
	issueMarketRegulationMap map[string]map[string]domain.IssueMarketRegulation
	operationStatusKabuMap   map[string]map[string]domain.OperationStatusKabu
	targetIssueCodes         []string
	targetIssueCodesMu       sync.RWMutex
	masterData               *domain.MasterData
}
