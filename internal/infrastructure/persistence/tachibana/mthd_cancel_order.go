package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CancelOrder は注文をキャンセルします。
func (tc *TachibanaClientImple) CancelOrder(ctx context.Context, orderID string) error {
	payload := map[string]string{
		"sCLMID":          clmidCancelOrder,
		"sOrderNumber":    orderID,
		"sEigyouDay":      "",
		"sSecondPassword": tc.Secret,
		"p_no":            tc.getPNo(),
		"p_sd_date":       formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.RequestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3) // 3回リトライ設定し送信
	if err != nil {
		return fmt.Errorf("cancel order failed: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		return fmt.Errorf("cancel order API returned an error: %s - %s", resultCode, response["sResultText"])
	}
	return nil
}
