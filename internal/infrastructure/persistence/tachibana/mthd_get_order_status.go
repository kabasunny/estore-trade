package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"time"
)

// GetOrderStatus は注文のステータスを取得
func (tc *TachibanaClientImple) GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error) {
	payload := map[string]string{
		"sCLMID":       clmidOrderListDetail,
		"sOrderNumber": orderID,
		"sEigyouDay":   "",
		"p_no":         tc.getPNo(),
		"p_sd_date":    formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.requestURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3) // 3回リトライ設定し送信
	if err != nil {
		return nil, fmt.Errorf("get order status failed: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		return nil, fmt.Errorf("order status API returned an error: %s - %s", resultCode, response["sResultText"])
	}

	order := &domain.Order{
		ID:     response["sOrderNumber"].(string),
		Status: response["sOrderStatus"].(string),
		// 他の必要なフィールドもマッピング
	}

	return order, nil
}
