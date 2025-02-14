// internal/infrastructure/persistence/tachibana/tachibana_client_impl.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"sync" // Mutex を使うため
	"time"

	"go.uber.org/zap"
)

type TachibanaClientIntImple struct {
	baseURL    *url.URL
	apiKey     string
	secret     string
	logger     *zap.Logger
	requestURL string       // キャッシュする仮想URL
	expiry     time.Time    // 仮想URLの有効期限
	mu         sync.RWMutex // 排他制御用
}

func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) TachibanaClient {
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err))
		return nil
	}
	return &TachibanaClientIntImple{
		baseURL: parsedURL,
		apiKey:  cfg.TachibanaAPIKey,
		secret:  cfg.TachibanaAPISecret,
		logger:  logger,
	}
}

// Login は API にログインし、仮想URLを返す。有効期限内ならキャッシュされたURLを返す
func (tc *TachibanaClientIntImple) Login(ctx context.Context, userID, password string) (string, error) {
	// --- ここから修正 ---
	// Read Lock: キャッシュされたURLが有効ならそれを返す
	tc.mu.RLock()
	if time.Now().Before(tc.expiry) && tc.requestURL != "" {
		tc.mu.RUnlock()
		return tc.requestURL, nil
	}
	tc.mu.RUnlock() // Unlockを確実に実行

	// Write Lock: 新しいURLを取得
	tc.mu.Lock()
	defer tc.mu.Unlock() // Unlockを確実に実行(defer)

	// 他のゴルーチンが既にURLを更新したかもしれないので、再度チェック
	if time.Now().Before(tc.expiry) && tc.requestURL != "" {
		return tc.requestURL, nil
	}
	// --- ここまで修正 ---

	payload := map[string]string{
		"sCLMID":    "CLMAuthLoginRequest",
		"sUserId":   userID,
		"sPassword": password,
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("Failed to marshal login payload", zap.Error(err))
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.baseURL.String()+"login", bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("Failed to create login request", zap.Error(err))
		return "", fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		tc.logger.Error("Failed to send login request", zap.Error(err))
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tc.logger.Error("Login failed: non-200 status code", zap.Int("status_code", resp.StatusCode))
		return "", fmt.Errorf("login failed: non-200 status code: %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tc.logger.Error("Failed to decode login response", zap.Error(err))
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	if response["sResultCode"] != "0" {
		tc.logger.Error("Login API returned an error", zap.String("result_code", response["sResultCode"]), zap.String("result_text", response["sResultText"]))
		return "", fmt.Errorf("login API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	requestURL, ok := response["sUrlRequest"]
	if !ok {
		tc.logger.Error("sUrlRequest not found in login response")
		return "", fmt.Errorf("sUrlRequest not found in login response")
	}

	// --- ここから修正 ---
	// キャッシュの更新 (有効期限は仮に1時間後とする)
	tc.requestURL = requestURL
	tc.expiry = time.Now().Add(1 * time.Hour) // 有効期限: 1時間後
	// --- ここまで修正 ---

	return requestURL, nil
}

// ---リファクタリング用 途中---
// 定数定義
const (
	clmidPlaceOrder            = "CLMKabuNewOrder"
	zyoutoekiKazeiCTokutei     = "1"  // 特定口座
	sizyouCToushou             = "00" // 東証
	baibaiKubunBuy             = "3"
	baibaiKubunSell            = "1"
	conditionSashine           = "0" // 指値
	genkinShinyouKubunGenbutsu = "0" // 現物
	orderExpireDay             = "0" // 当日限り
)

// ---リファクタリング用 途中---

// PlaceOrder は API に対して新しい株式注文を行う
func (tc *TachibanaClientIntImple) PlaceOrder(ctx context.Context, requestURL string, order *domain.Order) (*domain.Order, error) {
	//立花証券の注文APIの仕様に合わせてデータを作成

	// ---リファクタリング例---
	payload := map[string]interface{}{
		"sCLMID":              clmidPlaceOrder,                                                               // 定数
		"sZyoutoekiKazeiC":    zyoutoekiKazeiCTokutei,                                                        // 定数
		"sIssueCode":          order.Symbol,                                                                  // 銘柄コード
		"sSizyouC":            sizyouCToushou,                                                                // 定数
		"sBaibaiKubun":        map[string]string{"buy": baibaiKubunBuy, "sell": baibaiKubunSell}[order.Side], // 定数
		"sCondition":          conditionSashine,                                                              // 定数
		"sOrderPrice":         strconv.FormatFloat(order.Price, 'f', -1, 64),
		"sOrderSuryou":        strconv.Itoa(order.Quantity),
		"sGenkinShinyouKubun": genkinShinyouKubunGenbutsu, // 定数
		"sOrderExpireDay":     orderExpireDay,             // 定数
		// ... 他のフィールド ...
		"sSecondPassword": tc.secret, //第2パスワード
	}

	// ---リファクタリングここまで---
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("注文ペイロードのJSONエンコードに失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal order payload: %w", err)
	}

	// リクエストの送信
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("注文リクエストの作成に失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to create order request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		tc.logger.Error("注文リクエストの送信に失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to send order request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tc.logger.Error("注文APIが非200ステータスコードを返しました", zap.Int("status_code", resp.StatusCode))
		return nil, fmt.Errorf("order API returned non-200 status code: %d", resp.StatusCode)
	}

	var response map[string]interface{} //型をinterface{}に変更
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		tc.logger.Error("注文レスポンスのJSONデコードに失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to decode order response: %w", err)
	}
	//文字列型で受け取る
	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		//警告コードもログ出力
		warnCode, _ := response["sWarningCode"].(string) //警告コードも文字列
		warnText, _ := response["sWarningText"].(string)
		tc.logger.Error("注文APIがエラーを返しました", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)), zap.String("warning_code", warnCode), zap.String("warning_text", warnText))
		return nil, fmt.Errorf("order API returned an error: %s - %s", resultCode, response["sResultText"])
	}
	// 注文成功時の処理 (レスポンスから必要な情報を抽出)
	order.ID = response["sOrderNumber"].(string) // 注文番号
	order.Status = "pending"                     // ステータスを更新

	return order, nil
}

func (tc *TachibanaClientIntImple) GetOrderStatus(ctx context.Context, requestURL string, orderID string) (*domain.Order, error) {
	payload := map[string]string{
		"sCLMID":       "CLMOrderListDetail",
		"sOrderNumber": orderID,
		"sEigyouDay":   "", // 必要に応じて営業日を設定
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order status request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create order status request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send order status request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("order status API returned non-200 status code: %d", resp.StatusCode)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode order status response: %w", err)
	}

	if response["sResultCode"] != "0" {
		return nil, fmt.Errorf("order status API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	order := &domain.Order{
		ID:     response["sOrderNumber"].(string),
		Status: response["sOrderStatus"].(string), // APIのsOrderStatusを使用
		// 他の必要なフィールドもマッピング
	}

	return order, nil
}

func (tc *TachibanaClientIntImple) CancelOrder(ctx context.Context, requestURL string, orderID string) error {
	payload := map[string]string{
		"sCLMID":          "CLMKabuCancelOrder",
		"sOrderNumber":    orderID,
		"sEigyouDay":      "", // 必要に応じて営業日を設定
		"sSecondPassword": tc.secret,
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal cancel order request payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create cancel order request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send cancel order request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cancel order API returned non-200 status code: %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to decode cancel order response: %w", err)
	}

	if response["sResultCode"] != "0" {
		return fmt.Errorf("cancel order API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	return nil
}

// ConnectEventStream は、EVENT I/F への接続を確立し、受信したイベントをチャネルに流す
func (tc *TachibanaClientIntImple) ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error) {
	//  EventStream 構造体を使うように変更
	return nil, fmt.Errorf("ConnectEventStream method should be implemented in event_stream.go")
}
