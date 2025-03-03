// internal/infrastructure/persistence/tachibana/util_login.go
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
		"p_no":      tc.getPNo(), // 呼び出し元でインクリメントするため、ここでは取得のみ
		"p_sd_date": formatSDDate(time.Now()),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tc.logger.Error("Failed to marshal login payload", zap.Error(err))
		return false, fmt.Errorf("failed to marshal login payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.baseURL.String()+"login", bytes.NewBuffer(payloadJSON)) //BaseURL使用
	if err != nil {
		tc.logger.Error("Failed to create login request", zap.Error(err))
		return false, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3)
	if err != nil {
		return false, fmt.Errorf("login failed: %w", err)
	}

	if response["sResultCode"] != "0" {
		tc.logger.Error("Login API returned an error", zap.String("result_code", response["sResultCode"].(string)), zap.String("result_text", response["sResultText"].(string)))
		return false, fmt.Errorf("login API returned an error: %s - %s", response["sResultCode"], response["sResultText"])
	}

	requestURL, ok := response["sUrlRequest"]
	if !ok {
		tc.logger.Error("sUrlRequest not found in login response")
		return false, fmt.Errorf("sUrlRequest not found in login response")
	}
	masterURL, ok := response["sUrlMaster"]
	if !ok {
		tc.logger.Error("sUrlMaster not found in login response")
		return false, fmt.Errorf("sUrlMaster not found in login response")
	}
	priceURL, ok := response["sUrlPrice"]
	if !ok {
		tc.logger.Error("sUrlPrice not found in login response")
		return false, fmt.Errorf("sUrlPrice not found in login response")
	}
	eventURL, ok := response["sUrlEvent"]
	if !ok {
		tc.logger.Error("sUrlEvent not found in login response")
		return false, fmt.Errorf("sUrlEvent not found in login response")
	}

	// p_no はLogin APIのレスポンスで上書き
	if pNoStr, ok := response["p_no"].(string); ok {
		if pNo, err := strconv.ParseInt(pNoStr, 10, 64); err == nil {
			tc.pNo = pNo
		} else {
			tc.logger.Warn("Failed to parse p_no", zap.String("p_no", pNoStr), zap.Error(err))
		}
	} else {
		tc.logger.Warn("p_no not found or not a string in login response", zap.Any("response", response))
	}

	tc.requestURL = requestURL.(string)
	tc.masterURL = masterURL.(string)
	tc.priceURL = priceURL.(string)
	tc.eventURL = eventURL.(string)
	tc.expiry = time.Now().Add(1 * time.Hour)
	tc.loggined = true

	return true, nil
}
