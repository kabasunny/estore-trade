// internal/infrastructure/persistence/tachibana/mthd_get_positions.go
package tachibana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"estore-trade/internal/domain"
)

// GetPositions は建玉情報を取得します。
func (tc *TachibanaClientImple) GetPositions(ctx context.Context) ([]domain.Position, error) {
	if !tc.loggined {
		return nil, fmt.Errorf("not logged in")
	}
	requestURL, err := tc.GetRequestURL()
	if err != nil {
		return nil, fmt.Errorf("request URL not found, need to Login: %w", err)
	}

	// リクエストデータの作成
	payload := map[string]interface{}{ // string -> interface{}
		"sCLMID":     "CLMShinyouTategyokuList",
		"sIssueCode": "",
		"p_no":       tc.getPNo(),
		"p_sd_date":  formatSDDate(time.Now()),
		"sJsonOfmt":  "4", // JSON出力フォーマット (追加)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// URLエンコード (GETリクエスト)
	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL += "?" + encodedPayload

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	// リクエスト送信
	response, err := sendRequest(req, 3)
	if err != nil {
		return nil, fmt.Errorf("get positions failed: %w", err)
	}

	// レスポンスの処理 (sResultCode のチェック)
	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		resultText, _ := response["sResultText"].(string) // エラーメッセージも取得 (存在しない場合も考慮)
		return nil, fmt.Errorf("get positions API returned an error: %s - %s", resultCode, resultText)
	}

	// 建玉リストの取得
	positions := []domain.Position{}
	if positionsList, ok := response["aShinyouTategyokuList"].([]interface{}); ok {
		for _, p := range positionsList {
			positionMap, ok := p.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid position data format")
			}

			position, err := convertToPosition(positionMap)
			if err != nil {
				return nil, err
			}
			positions = append(positions, *position)
		}
	}

	return positions, nil
}
