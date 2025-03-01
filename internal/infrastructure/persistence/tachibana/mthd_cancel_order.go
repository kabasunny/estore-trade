// internal/infrastructure/persistence/tachibana/mthd_cancel_order.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// CancelOrder は注文をキャンセルします。
func (tc *TachibanaClientImple) CancelOrder(ctx context.Context, orderID string) error {
	payload := map[string]string{
		"sCLMID":          clmidCancelOrder,
		"sOrderNumber":    orderID,
		"sEigyouDay":      "", // 空で良いか立花証券のAPI仕様を確認
		"sSecondPassword": tc.secret,
		"p_no":            tc.getPNo(),
		"p_sd_date":       formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	requestURL, err := url.JoinPath(tc.requestURL, "cancel") //tc.RequestURLが正しい前提
	if err != nil {
		return fmt.Errorf("failed to create request URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestURL, bytes.NewBuffer(payloadJSON))
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
		// 警告コードがある場合もログに出力 (PlaceOrder, GetOrderStatus に倣う)
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)

		// 型アサーションを使って string に変換
		resultText, ok := response["sResultText"].(string)
		if !ok {
			// 型アサーションに失敗した場合の処理 (例えば、ログに出力してエラーを返す)
			tc.logger.Error("sResultText is not a string", zap.Any("sResultText", response["sResultText"]))
			return fmt.Errorf("sResultText is not a string in the response") // エラーだけを返す
		}

		tc.logger.Error("注文キャンセルAPIがエラーを返しました",
			zap.String("result_code", resultCode),
			zap.String("result_text", resultText), // resultText を使用
			zap.String("order_id", orderID),       // orderIDもログに出力
			zap.String("warning_code", warnCode),
			zap.String("warning_text", warnText),
		)
		return fmt.Errorf("cancel order API returned an error: %s - %s", resultCode, response["sResultText"])
	}
	return nil
}
