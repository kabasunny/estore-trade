package tachibana

import (
	"net/url"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TachibanaClientImple 構造体の定義
type TachibanaClientImple struct {
	baseURL    *url.URL
	apiKey     string
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

	// マスタデータ保持用 (必要最低限に絞り込み)
	systemStatus             SystemStatus
	dateInfo                 DateInfo
	callPriceMap             map[string]CallPrice
	issueMap                 map[string]IssueMaster
	issueMarketMap           map[string]map[string]IssueMarketMaster
	issueMarketRegulationMap map[string]map[string]IssueMarketRegulation
	operationStatusKabuMap   map[string]map[string]OperationStatusKabu
	targetIssueCodes         []string
	targetIssueCodesMu       sync.RWMutex
}
