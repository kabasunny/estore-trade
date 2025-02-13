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
	"time"

	"go.uber.org/zap"
)

type TachibanaClientIntImple struct {
	baseURL *url.URL //URL型に変更
	apiKey  string
	secret  string
	logger  *zap.Logger
}

func NewTachibanaClient(cfg *config.Config, logger *zap.Logger) TachibanaClient {
	// URLのパースとエラーハンドリング
	parsedURL, err := url.Parse(cfg.TachibanaBaseURL)
	if err != nil {
		// URLのパースに失敗した場合、致命的なエラーとして扱う
		logger.Fatal("Invalid Tachibana API base URL", zap.Error(err)) //loggerでエラー記録
		return nil                                                     //nilを返して呼び出し元で処理
	}
	return &TachibanaClientIntImple{
		baseURL: parsedURL, //パースされたURL
		apiKey:  cfg.TachibanaAPIKey,
		secret:  cfg.TachibanaAPISecret,
		logger:  logger, //loggerを受け取る
	}
}

// APIに対してログインし、ユーザーIDとパスワードを使用して必要な認証情報を取得し、成功した場合、APIとやり取りするためのリクエストURLを返す
func (tc *TachibanaClientIntImple) Login(ctx context.Context, userID, password string) (string, error) {
	// 謎の文字列キーは、API仕様書にて参照　https://www.e-shiten.jp/e_api/mfds_json_api_refference.html

	// リクエストデータの作成
	payload := map[string]string{
		"sCLMID":    "CLMAuthLoginRequest", // 機能ID：ログインリクエスト
		"sUserId":   userID,                // ユーザーID
		"sPassword": password,              // パスワード
	}

	// JSON 形式にエンコード
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("Failed to marshal login payload", zap.Error(err)) // エラーログ
		return "", fmt.Errorf("failed to marshal login payload: %w", err)
	}

	// リクエストの作成と送信
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.baseURL.String()+"login", bytes.NewBuffer(payloadJSON)) // baseURLを文字列に変換して使用
	if err != nil {
		tc.logger.Error("Failed to create login request", zap.Error(err))
		return "", fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // ヘッダーにContent-Typeを設定
	client := &http.Client{Timeout: 10 * time.Second}  // タイムアウトを設定
	resp, err := client.Do(req)
	if err != nil {
		tc.logger.Error("Failed to send login request", zap.Error(err))
		return "", fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close() // 関数の終了時にレスポンスボディを閉じる

	// レスポンスの処理
	if resp.StatusCode != http.StatusOK {
		tc.logger.Error("Login failed: non-200 status code", zap.Int("status_code", resp.StatusCode))
		return "", fmt.Errorf("login failed: non-200 status code: %d", resp.StatusCode)
	}

	// レスポンスボディをデコード
	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil { // データストリームをリアルタイムで直接Goのデータ構造にデコード
		tc.logger.Error("Failed to decode login response", zap.Error(err))
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	// 結果コードの確認：0は正常なとき
	if response["sResultCode"] != "0" { // 0は正常、他はエラー
		tc.logger.Error("Login API returned an error", zap.String("result_code", response["sResultCode"]), zap.String("result_text", response["sResultText"]))
		return "", fmt.Errorf("login API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	// 仮想URLを取得
	requestURL, ok := response["sUrlRequest"]
	if !ok {
		tc.logger.Error("sUrlRequest not found in login response") // 仮想URLがレスポンスに含まれていない場合のエラーログ
		return "", fmt.Errorf("sUrlRequest not found in login response")
	}
	return requestURL, nil // 正常終了時にリクエストURLを返す
}

// APIに対して新しい株式注文
func (tc *TachibanaClientIntImple) PlaceOrder(ctx context.Context, requestURL string, order *domain.Order) (*domain.Order, error) {
	//立花証券の注文APIの仕様に合わせてデータを作成
	payload := map[string]interface{}{ //interface{}で異なる型を許容
		"sCLMID":                    "CLMKabuNewOrder",
		"sZyoutoekiKazeiC":          "1", // 例: 特定口座
		"sIssueCode":                order.Symbol,
		"sSizyouC":                  "00", // 例: 東証
		"sBaibaiKubun":              map[string]string{"buy": "3", "sell": "1"}[order.Side],
		"sCondition":                "0",                                           // 例: 指値
		"sOrderPrice":               strconv.FormatFloat(order.Price, 'f', -1, 64), // 文字列に変換
		"sOrderSuryou":              strconv.Itoa(order.Quantity),                  // 文字列に変換
		"sGenkinShinyouKubun":       "0",                                           // 例: 現物
		"sOrderExpireDay":           "0",
		"sGyakusasiOrderType":       "0",
		"sGyakusasiZyouken":         "0",
		"sGyakusasiPrice":           "*",
		"sTatebiType":               "*",
		"sTategyokuZyoutoekiKazeiC": "*",
		"sSecondPassword":           tc.secret, //第2パスワード
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("注文ペイロードのJSONエンコードに失敗", zap.Error(err))
		return nil, fmt.Errorf("failed to marshal order payload: %w", err)
	}

	// リクエストの送信
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(payloadJSON)) //ctx context.Context
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

// 注文IDに基づいて、注文のステータスを取得
func (tc *TachibanaClientIntImple) GetOrderStatus(ctx context.Context, requestURL string, orderID string) (*domain.Order, error) {
	// 1. リクエストデータの準備 (CLMOrderListDetail を使用)
	payload := map[string]string{
		"sCLMID":       "CLMOrderListDetail",
		"sOrderNumber": orderID,
		"sEigyouDay":   "", // 必要に応じて営業日を設定
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order status request payload: %w", err)
	}

	// 2. リクエストの送信
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

	// 3. レスポンスの処理
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

	// 4. レスポンスから必要な情報を抽出して、domain.Orderオブジェクトにマッピング
	order := &domain.Order{
		ID:     response["sOrderNumber"].(string),
		Status: response["sOrderStatus"].(string), // APIのsOrderStatusを使用
		// 他の必要なフィールドもマッピング
	}

	return order, nil
}

// 注文IDに基づいて、注文のキャンセル
func (tc *TachibanaClientIntImple) CancelOrder(ctx context.Context, requestURL string, orderID string) error {
	// 1. リクエストデータの準備 (CLMKabuCancelOrder を使用)
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

	// 2. リクエストの送信
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

	// 3. レスポンスの処理
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

	return nil // キャンセル成功
}
