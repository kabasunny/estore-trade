package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// PlaceOrder は API に対して新しい株式注文を行う
func (tc *TachibanaClientImple) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	payload := map[string]interface{}{
		"sCLMID":              clmidPlaceOrder,
		"sZyoutoekiKazeiC":    zyoutoekiKazeiCTokutei,
		"sIssueCode":          order.Symbol,
		"sSizyouC":            sizyouCToushou,
		"sBaibaiKubun":        map[string]string{"buy": baibaiKubunBuy, "sell": baibaiKubunSell}[order.Side],
		"sCondition":          conditionSashine,
		"sOrderPrice":         strconv.FormatFloat(order.Price, 'f', -1, 64),
		"sOrderSuryou":        strconv.Itoa(order.Quantity),
		"sGenkinShinyouKubun": genkinShinyouKubunGenbutsu,
		"sOrderExpireDay":     orderExpireDay,
		"sSecondPassword":     tc.Secret,
		"p_no":                tc.getPNo(),
		"p_sd_date":           formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tc.RequestURL, bytes.NewBuffer(payloadJSON))
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

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)
		tc.Logger.Error("注文APIがエラーを返しました", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)), zap.String("warning_code", warnCode), zap.String("warning_text", warnText))
		return nil, fmt.Errorf("order API returned an error: %s - %s", resultCode, response["sResultText"])
	}

	order.ID = response["sOrderNumber"].(string)
	order.Status = "pending"

	return order, nil
}
