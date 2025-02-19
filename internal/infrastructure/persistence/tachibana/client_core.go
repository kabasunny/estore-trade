// internal/infrastructure/persistence/tachibana/client_core.go
package tachibana

import (
	"context"
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"estore-trade/internal/config"
	"estore-trade/internal/domain"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	_ "golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	_ "golang.org/x/text/transform"
)

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
	callPriceMap             map[string]CallPrice                    // 呼値 (Key: sYobineTaniNumber)
	issueMap                 map[string]IssueMaster                  // 銘柄マスタ（株式）(Key: 銘柄コード)
	issueMarketMap           map[string]map[string]IssueMarketMaster // 株式銘柄市場マスタ (Key1: 銘柄コード, Key2: 上場市場)
	issueMarketRegulationMap map[string]map[string]IssueMarketRegulation
	operationStatusKabuMap   map[string]map[string]OperationStatusKabu // 運用ステータス（株）(Key1: 上場市場, Key2: 運用単位)
	targetIssueCodes         []string                                  // ターゲット銘柄コード
	targetIssueCodesMu       sync.RWMutex                              // 排他制御
}

func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) TachibanaClient {
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err))
		return nil
	}
	return &TachibanaClientImple{
		baseURL:                  parsedURL,
		apiKey:                   cfg.TachibanaAPIKey,
		secret:                   cfg.TachibanaAPISecret,
		logger:                   logger,
		loggined:                 false,                                             // 初期値はfalse
		pNo:                      0,                                                 // 初期値は0
		callPriceMap:             make(map[string]CallPrice),                        // 追加: 呼値
		issueMap:                 make(map[string]IssueMaster),                      // 追加: 銘柄
		issueMarketMap:           make(map[string]map[string]IssueMarketMaster),     // 株式銘柄市場マスタ
		issueMarketRegulationMap: make(map[string]map[string]IssueMarketRegulation), //株式銘柄別・市場別規制
		operationStatusKabuMap:   make(map[string]map[string]OperationStatusKabu),   // 運用ステータス（株）
		targetIssueCodes:         make([]string, 0),                                 // 初期化
	}
}

// Login は API にログインし、仮想URLを返す。有効期限内ならキャッシュされたURLを返す
func (tc *TachibanaClientImple) Login(ctx context.Context, cfg interface{}) error {
	userID := cfg.(*config.Config).TachibanaUserID //型アサーション
	password := cfg.(*config.Config).TachibanaPassword

	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" && tc.masterURL != "" && tc.priceURL != "" && tc.eventURL != "" {
		tc.mu.RUnlock()
		return nil
	}
	tc.mu.RUnlock()

	// Write Lock: 新しいURLを取得
	tc.mu.Lock()

	loggined, err := login(ctx, tc, userID, password)
	tc.loggined = loggined

	defer tc.mu.Unlock()

	return err
}

// getPNo は p_no を取得し、インクリメントする (スレッドセーフ)
func (tc *TachibanaClientImple) getPNo() string {
	tc.pNoMu.Lock()
	defer tc.pNoMu.Unlock()
	tc.pNo++
	return strconv.FormatInt(tc.pNo, 10)
}

// ConnectEventStream は、EVENT I/F への接続を確立し、受信したイベントをチャネルに流す
func (tc *TachibanaClientImple) ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error) {
	//  EventStream 構造体を使うように変更
	return nil, fmt.Errorf("ConnectEventStream method should be implemented in event_stream.go")
}

// SetTargetIssues は、指定された銘柄コードのみを対象とするようにマスタデータをフィルタリングする
func (tc *TachibanaClientImple) SetTargetIssues(ctx context.Context, issueCodes []string) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// issueMap のフィルタリング
	for issueCode := range tc.issueMap {
		if !contains(issueCodes, issueCode) { // ヘルパー関数を使用
			delete(tc.issueMap, issueCode)
		}
	}

	// issueMarketMap, issueMarketRegulationMap のフィルタリング (issueMap と同様)
	for issueCode := range tc.issueMarketMap {
		if !contains(issueCodes, issueCode) {
			delete(tc.issueMarketMap, issueCode)
			continue // 銘柄コードが削除されたら、その下の市場情報も不要
		}
		// (特定の市場だけが必要な場合は、ここでさらにフィルタリング)
	}

	// issueMarketRegulationMap のフィルタリング
	for issueCode := range tc.issueMarketRegulationMap {
		if !contains(issueCodes, issueCode) {
			delete(tc.issueMarketRegulationMap, issueCode)
			continue // 銘柄コードが削除されたら、その下の市場情報も不要
		}
		// (特定の市場だけが必要な場合は、ここでさらにフィルタリング)
	}

	tc.targetIssueCodesMu.Lock() // 排他制御
	tc.targetIssueCodes = issueCodes
	tc.targetIssueCodesMu.Unlock()
	return nil
}

// (必要に応じて) GetIssueMaster, GetIssueMarketMaster, GetIssueMarketRegulation の修正:
func (tc *TachibanaClientImple) GetIssueMaster(issueCode string) (IssueMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) { // ヘルパー関数を使用
			tc.targetIssueCodesMu.RUnlock()
			return IssueMaster{}, false // ターゲット銘柄でなければエラー
		}

	}
	tc.targetIssueCodesMu.RUnlock()

	issue, ok := tc.issueMap[issueCode]
	return issue, ok
}

// GetIssueMarketMaster, GetIssueMarketRegulation も同様に修正
func (tc *TachibanaClientImple) GetIssueMarketMaster(issueCode, marketCode string) (IssueMarketMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			tc.targetIssueCodesMu.RUnlock()
			return IssueMarketMaster{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.targetIssueCodesMu.RUnlock()

	marketMap, ok := tc.issueMarketMap[issueCode]
	if !ok {
		return IssueMarketMaster{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}

func (tc *TachibanaClientImple) GetIssueMarketRegulation(issueCode, marketCode string) (IssueMarketRegulation, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			tc.targetIssueCodesMu.RUnlock()
			return IssueMarketRegulation{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.targetIssueCodesMu.RUnlock()

	marketMap, ok := tc.issueMarketRegulationMap[issueCode]
	if !ok {
		return IssueMarketRegulation{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}

// GetOperationStatusKabu　も同様
func (tc *TachibanaClientImple) GetOperationStatusKabu(listedMarket string, unit string) (OperationStatusKabu, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	marketMap, ok := tc.operationStatusKabuMap[listedMarket]
	if !ok {
		return OperationStatusKabu{}, false
	}
	status, ok := marketMap[unit]
	return status, ok
}

// CheckPriceIsValid は、呼値があっているか確認
func (tc *TachibanaClientImple) CheckPriceIsValid(issueCode string, price float64, isNextDay bool) (bool, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// 銘柄の呼値単位番号を取得 (翌営業日の場合は Yoku を使う)
	issueMarket, ok := tc.issueMarketMap[issueCode]["00"] // 例として市場コード"00" (東証) を使用
	if !ok {
		return false, fmt.Errorf("IssueMarketMaster not found for issueCode: %s", issueCode)
	}
	unitNumberStr := issueMarket.CallPriceUnitNumber
	if isNextDay {
		unitNumberStr = issueMarket.CallPriceUnitNumberYoku
	}
	if unitNumberStr == "" {
		// 呼値情報がない場合はチェック不要 (またはエラーとする)
		return true, nil // または return false, fmt.Errorf(...)
	}

	unitNumber, err := strconv.Atoi(unitNumberStr)
	if err != nil {
		return false, fmt.Errorf("invalid CallPriceUnitNumber: %s", unitNumberStr)
	}

	callPrice, ok := tc.callPriceMap[strconv.Itoa(unitNumber)]
	if !ok {
		return false, fmt.Errorf("CallPrice not found for unitNumber: %d", unitNumber)
	}

	// isValidPrice 関数を使ってチェック
	return isValidPrice(price, callPrice), nil
}

// sendRequest は、HTTPリクエストを送信し、レスポンスをデコードする (リトライ処理付き)
func sendRequest(ctx context.Context, tc *TachibanaClientImple, req *http.Request) (map[string]interface{}, error) {
	// リトライ処理を retryDo 関数に委譲
	var response map[string]interface{}
	retryFunc := func() (*http.Response, error) {
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return resp, err // エラーをそのまま返す
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return resp, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode) // ステータスコードが200以外の場合もエラーとして返す
		}

		// レスポンスのデコード処理もここで行う
		reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
		if err := json.NewDecoder(reader).Decode(&response); err != nil {
			resp.Body.Close() // デコードに失敗した場合もクローズ
			return resp, fmt.Errorf("failed to decode response: %w", err)
		}
		return resp, nil // 成功時はここ
	}

	resp, err := retryDo(retryFunc, 2, 2*time.Second) // 最大3回、初期遅延2秒
	if err != nil {
		return nil, err // retryDo でエラー処理済み
	}
	defer resp.Body.Close()

	return response, nil
}

func (tc *TachibanaClientImple) GetRequestURL() (string, error) {
	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" {
		return tc.requestURL, nil
	}
	return "", fmt.Errorf("request URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetMasterURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.masterURL != "" {
		return tc.masterURL, nil
	}
	return "", fmt.Errorf("master URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetPriceURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.priceURL != "" {
		return tc.priceURL, nil
	}

	return "", fmt.Errorf("price URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetEventURL() (string, error) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.eventURL != "" {
		return tc.eventURL, nil
	}
	return "", fmt.Errorf("event URL not found, neead to Login")
}

// マスタデータへのアクセス用メソッド (Getter)
func (tc *TachibanaClientImple) GetSystemStatus() SystemStatus {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.systemStatus
}

func (tc *TachibanaClientImple) GetDateInfo() DateInfo {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return tc.dateInfo
}

func (tc *TachibanaClientImple) GetCallPrice(unitNumber string) (CallPrice, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	callPrice, ok := tc.callPriceMap[unitNumber]
	return callPrice, ok
}
