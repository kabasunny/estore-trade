// internal/infrastructure/persistence/tachibana/client_login.go

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
)

func login(ctx context.Context, tc *TachibanaClientImple, userID, password string) (bool, error) {
	payload := map[string]string{
		"sCLMID":    clmidLogin,
		"sUserId":   userID,
		"sPassword": password,
		"p_no":      tc.getPNo(),              // p_no を設定 v4.6以降不要か
		"p_sd_date": formatSDDate(time.Now()), // p_sd_date を設定 v4.6以降不要か
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("Failed to marshal login payload", zap.Error(err))
		return false, fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.baseURL.String()+"login", bytes.NewBuffer(payloadJSON))
	if err != nil {
		tc.logger.Error("Failed to create login request", zap.Error(err))
		return false, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// コンテキストとタイムアウトの設定
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	// sendRequestを呼び出す（リトライ処理はsendRequest内で実施）
	response, err := sendRequest(req) // reqを渡すように変更
	if err != nil {
		return false, fmt.Errorf("login failed: %w", err)
	}

	// ... (以降は同じ) ...
	if response["sResultCode"] != "0" {
		tc.logger.Error("Login API returned an error", zap.String("result_code", response["sResultCode"].(string)), zap.String("result_text", response["sResultText"].(string))) //.(string)を追加
		return false, fmt.Errorf("login API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	requestURL, ok := response["sUrlRequest"] // sUrlRequest	仮想URL（REQUEST)	業務機能　　（REQUEST I/F）仮想URL
	if !ok {
		tc.logger.Error("sUrlRequest not found in login response")
		return false, fmt.Errorf("sUrlRequest not found in login response")
	}

	masterURL, ok := response["sUrlMaster"] // sUrlMaster	仮想URL（Master)	マスタ機能　（REQUEST I/F）仮想URL
	if !ok {
		tc.logger.Error("sUrlMaster not found in login response")
		return false, fmt.Errorf("sUrlMaster not found in login response")
	}

	priceURL, ok := response["sUrlPrice"] // sUrlPrice	仮想URL（Price)	時価情報機能（REQUEST I/F）仮想URL
	if !ok {
		tc.logger.Error("sUrlPrice not found in login response")
		return false, fmt.Errorf("sUrlPrice not found in login response")
	}

	eventURL, ok := response["sUrlEvent"] // sUrlEvent	仮想URL（EVENT)	注文約定通知（EVENT I/F）仮想URL
	if !ok {
		tc.logger.Error("sUrlEvent not found in login response")
		return false, fmt.Errorf("sUrlEvent not found in login response")
	}

	// p_no の初期値を設定 (Login成功時のみ)
	if pNoStr, ok := response["p_no"].(string); ok { // 型アサーションを追加
		if pNo, err := strconv.ParseInt(pNoStr, 10, 64); err == nil {
			tc.pNo = pNo
		} else {
			tc.logger.Warn("Failed to parse p_no", zap.String("p_no", pNoStr), zap.Error(err)) // パースエラーをログに記録
			// 必要に応じてデフォルト値を設定するか、エラー処理を行う
		}
	} else {
		tc.logger.Warn("p_no not found or not a string in login response", zap.Any("response", response)) // p_no が存在しない、または文字列でない場合
		// 必要に応じてデフォルト値を設定するか、エラー処理を行う
	}

	// キャッシュの更新 (有効期限は仮に1時間後とする)
	tc.requestURL = requestURL.(string) // interface{}型からstring型に
	tc.masterURL = masterURL.(string)
	tc.priceURL = priceURL.(string)
	tc.eventURL = eventURL.(string)

	tc.expiry = time.Now().Add(1 * time.Hour) // 有効期限: 1時間後

	return true, nil
}
