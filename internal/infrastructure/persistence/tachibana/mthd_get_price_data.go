// internal/infrastructure/persistence/tachibana/mthd_get_price_data.go
package tachibana

import (
	"bytes"
	"context"
	"encoding/json"
	"estore-trade/internal/domain" //追加
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// GetPriceData は、指定された銘柄コードのリストに対して、時価情報（当日終値、出来高など）を取得します。
func (tc *TachibanaClientImple) GetPriceData(ctx context.Context, issueCodes []string) ([]domain.PriceData, error) { // 戻り値の型を変更
	priceURL, err := tc.GetPriceURL()
	if err != nil {
		return nil, fmt.Errorf("failed to get price URL: %w", err)
	}
	issueCodeStr := strings.Join(issueCodes, ",")

	payload := map[string]interface{}{
		"sCLMID":     clmFdsGetMarketPrice,
		"sIssueCode": issueCodeStr,
		"p_no":       tc.getPNo(),
		"p_sd_date":  formatSDDate(time.Now()),
	}
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, priceURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req, cancel := withContextAndTimeout(req, 60*time.Second) // タイムアウトは適宜調整
	defer cancel()

	response, err := sendRequest(req, 3) // リトライ処理 (sendRequest内で実装)
	if err != nil {
		return nil, fmt.Errorf("failed to get price data: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		warnCode, _ := response["sWarningCode"].(string)
		warnText, _ := response["sWarningText"].(string)
		tc.logger.Error("price data API returned an error", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)), zap.String("warning_code", warnCode), zap.String("warning_text", warnText))
		return nil, fmt.Errorf("price data API returned an error: %s - %s", resultCode, response["sResultText"])
	}

	// レスポンスデータの処理 (ここでは必要なフィールドのみ抽出)
	var priceDataList []domain.PriceData // 変更
	// APIレスポンスの構造に合わせて修正
	if data, ok := response["data"].([]interface{}); ok {
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				var priceData domain.PriceData // 変更
				if err := mapToStruct(itemMap, &priceData); err != nil {
					return nil, fmt.Errorf("failed to map PriceData: %w", err)
				}
				// sIssueCode がレスポンスにない場合は、リクエスト時の issueCodes から設定
				for _, code := range issueCodes {
					if strings.Contains(issueCodeStr, code) { //リクエストした銘柄コードと一致したら
						priceData.IssueCode = code
					}
				}
				priceDataList = append(priceDataList, priceData)
			}
		}
	}
	return priceDataList, nil
}
