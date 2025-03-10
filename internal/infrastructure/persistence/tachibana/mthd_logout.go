// internal/infrastructure/persistence/tachibana/mthd_logout.go
package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// Logout はAPIからログアウトします。
func (tc *TachibanaClientImple) Logout(ctx context.Context) error {
	tc.mu.Lock()
	defer tc.mu.Unlock()

	if !tc.loggined {
		return nil
	}

	// requestURL が空の場合はエラーを返す
	if tc.requestURL == "" {
		return fmt.Errorf("requestURL が空です") // 日本語のエラーメッセージ
	}

	payload := map[string]string{
		"sCLMID":    clmidLogoutRequest,
		"p_no":      tc.getPNo(),
		"p_sd_date": formatSDDate(time.Now()),
		"sJsonOfmt": "4",
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("JSON encode error: %w", err)
	}

	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL := tc.requestURL + "?" + encodedPayload

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return fmt.Errorf("request creation error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3)
	if err != nil {
		return fmt.Errorf("ログアウトに失敗しました: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); !ok {
		return fmt.Errorf("API error: sResultCode not found in response")
	} else if resultCode != "0" {
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)

		tc.logger.Error("API returned an error",
			zap.String("result_code", resultCode),
			zap.String("result_text", response["sResultText"].(string)),
			zap.String("warning_code", warnCode),
			zap.String("warning_text", warnText),
		)
		return fmt.Errorf("API error: %s - %s", resultCode, response["sResultText"])
	}

	tc.loggined = false
	tc.requestURL = ""
	tc.masterURL = ""
	tc.priceURL = ""
	tc.eventURL = ""
	tc.expiry = time.Time{}

	return nil
}
