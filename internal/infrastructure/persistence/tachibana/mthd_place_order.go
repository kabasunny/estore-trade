// internal/infrastructure/persistence/tachibana/mthd_place_order.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// PlaceOrder は API に対して新しい株式注文を行う
func (tc *TachibanaClientImple) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {

	payload, err := ConvertOrderToPlaceOrderPayload(order, tc) // 変換関数を呼び出す
	if err != nil {
		return nil, fmt.Errorf("failed to convert order to payload: %w", err)
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
		return nil, fmt.Errorf("place order failed: %w", err)
	}

	// internal/infrastructure/persistence/tachibana/mthd_place_order.go

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)
		tc.logger.Error("注文APIがエラーを返しました", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)), zap.String("warning_code", warnCode), zap.String("warning_text", warnText))
		return nil, fmt.Errorf("order API returned an error: %s - %s", resultCode, response["sResultText"])
	}

	//order.ID = response["sOrderNumber"].(string) //エラー箇所
	// order.Status = "pending" //初期状態ではpending

	//追加
	orderID, ok := response["sOrderNumber"].(string)
	if !ok {
		return nil, errors.New("order number not found in response")
	}
	order.ID = orderID
	order.Status = "pending" //初期状態ではpending

	return order, nil
}
