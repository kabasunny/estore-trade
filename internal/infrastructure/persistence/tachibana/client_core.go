// internal/infrastructure/persistence/tachibana/client_core.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
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
	"golang.org/x/text/transform"
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
	// マスタデータ保持用
	systemStatus SystemStatus
	dateInfo     DateInfo
	callPriceMap map[string]CallPrice   // 呼値 (Key: sYobineTaniNumber)
	issueMap     map[string]IssueMaster // 銘柄マスタ（株式）(Key: 銘柄コード)
}

func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) TachibanaClient {
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err))
		return nil
	}
	return &TachibanaClientImple{
		baseURL:      parsedURL,
		apiKey:       cfg.TachibanaAPIKey,
		secret:       cfg.TachibanaAPISecret,
		logger:       logger,
		loggined:     false,                        // 初期値はfalse
		pNo:          0,                            // 初期値は0
		callPriceMap: make(map[string]CallPrice),   // 追加: 呼値
		issueMap:     make(map[string]IssueMaster), // 追加: 銘柄
	}
}

// Login は API にログインし、仮想URLを返す。有効期限内ならキャッシュされたURLを返す
func (tc *TachibanaClientImple) Login(ctx context.Context, cfg *config.Config) (bool, error) {
	userID := cfg.TachibanaUserID
	password := cfg.TachibanaPassword

	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" && tc.masterURL != "" && tc.priceURL != "" && tc.eventURL != "" {
		tc.mu.RUnlock()
		return tc.loggined, nil
	}
	tc.mu.RUnlock()

	// Write Lock: 新しいURLを取得
	tc.mu.Lock()
	defer tc.mu.Unlock()

	return login(ctx, tc, userID, password)
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

func sendRequest(ctx context.Context, tc *TachibanaClientImple, payload interface{}) (map[string]interface{}, error) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("ペイロードのJSONエンコードに失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("リクエストの作成に失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()
	client := &http.Client{Timeout: 60 * time.Second} // タイムアウト設定
	resp, err := client.Do(req)
	if err != nil {
		tc.logger.Warn("Request failed", zap.Error(err))
		return nil, fmt.Errorf("request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		tc.logger.Error("API returned non-200 status code", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
	var response map[string]interface{}
	if err := json.NewDecoder(reader).Decode(&response); err != nil {
		tc.logger.Error("レスポンスのJSONデコードに失敗", zap.Error(err))
		resp.Body.Close()
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	resp.Body.Close()
	return response, nil
}

func (tc *TachibanaClientImple) GetRequestURL() (string, error) {
	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.requestURL != "" {
		tc.mu.RUnlock()
		return tc.requestURL, nil
	}
	tc.mu.RUnlock()
	return "", fmt.Errorf("request URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetMasterURL() (string, error) {
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.masterURL != "" {
		tc.mu.RUnlock()
		return tc.masterURL, nil
	}
	tc.mu.RUnlock()
	return "", fmt.Errorf("master URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetPriceURL() (string, error) {
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.priceURL != "" {
		tc.mu.RUnlock()
		return tc.priceURL, nil
	}
	tc.mu.RUnlock()
	return "", fmt.Errorf("price URL not found, neead to Login")
}

func (tc *TachibanaClientImple) GetEventURL() (string, error) {
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.loggined && tc.eventURL != "" {
		tc.mu.RUnlock()
		return tc.eventURL, nil
	}
	tc.mu.RUnlock()
	return "", fmt.Errorf("event URL not found, neead to Login")
}
