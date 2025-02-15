// internal/infrastructure/persistence/tachibana/login.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func login(ctx context.Context, tc *TachibanaClientIntImple, userID, password string) (string, error) {
	// リトライ処理 (最大3回、間隔は2秒から開始して指数関数的に増加)
	var resp *http.Response
	var err error
	for retries := 0; retries < 3; retries++ {

		payload := map[string]string{
			"sCLMID":    clmidLogin,
			"sUserId":   userID,
			"sPassword": password,
			"p_no":      tc.getPNo(),              // p_no を設定
			"p_sd_date": formatSDDate(time.Now()), // p_sd_date を設定
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
		// コンテキストとタイムアウトを設定
		req = withContextAndTimeout(req, 60*time.Second)
		client := &http.Client{} //  client := &http.Client{Timeout: 60 * time.Second}
		resp, err = client.Do(req)

		if err != nil {
			tc.logger.Warn("Login request failed, retrying...", zap.Error(err), zap.Int("retry", retries+1))
			time.Sleep(time.Duration(1+retries*2) * time.Second) // 指数バックオフ
			continue
		}

		// HTTPステータスコードのチェック
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close() // StatusCode が OK 意外の場合は Body を ক্লো졍
			tc.logger.Error("Login failed: non-200 status code", zap.Int("status_code", resp.StatusCode))
			if resp.StatusCode >= 500 {
				time.Sleep(time.Duration(1+retries*2) * time.Second)
				continue
			}
			return "", fmt.Errorf("login failed: non-200 status code: %d", resp.StatusCode)
		}

		// レスポンスボディをShift-JISからUTF-8に変換
		reader := transform.NewReader(resp.Body, japanese.ShiftJIS.NewDecoder())
		var response map[string]string
		if err := json.NewDecoder(reader).Decode(&response); err != nil {
			tc.logger.Error("Failed to decode login response", zap.Error(err))
			resp.Body.Close() // Body を ক্লো졍
			return "", fmt.Errorf("failed to decode login response: %w", err)
		}
		resp.Body.Close() // Body を ক্লো졍

		if response["sResultCode"] != "0" {
			tc.logger.Error("Login API returned an error", zap.String("result_code", response["sResultCode"]), zap.String("result_text", response["sResultText"]))
			return "", fmt.Errorf("login API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
		}

		requestURL, ok := response["sUrlRequest"]
		if !ok {
			tc.logger.Error("sUrlRequest not found in login response")
			return "", fmt.Errorf("sUrlRequest not found in login response")
		}

		// p_no の初期値を設定 (Login成功時のみ)
		if pNoStr, ok := response["p_no"]; ok {
			if pNo, err := strconv.ParseInt(pNoStr, 10, 64); err == nil {
				tc.pNo = pNo
			}
		}

		// キャッシュの更新 (有効期限は仮に1時間後とする)
		tc.requestURL = requestURL
		tc.expiry = time.Now().Add(1 * time.Hour) // 有効期限: 1時間後

		return requestURL, nil
	}
	return "", fmt.Errorf("login failed after multiple retries: %w", err) // 最終的なエラー
}
