// internal/infrastructure/persistence/tachibana/order.go
package tachibana

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// PlaceOrder は API に対して新しい株式注文を行う
func (tc *TachibanaClientImple) PlaceOrder(ctx context.Context, requestURL string, order *domain.Order) (*domain.Order, error) {
	// リトライ処理
	var err error
	for retries := 0; retries < 3; retries++ {
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
			"sSecondPassword":     tc.secret,
			"p_no":                tc.getPNo(),              // p_no を設定
			"p_sd_date":           formatSDDate(time.Now()), // p_sd_date を設定
		}

		response, err := sendRequest(ctx, tc, requestURL, payload)
		if err != nil {
			if isRetryableError(err) && retries < 2 {
				tc.logger.Warn("PlaceOrder request failed, retrying...", zap.Error(err), zap.Int("retry", retries+1))
				time.Sleep(time.Duration(1+retries*2) * time.Second) // 指数バックオフ
				continue
			}
			return nil, fmt.Errorf("place order failed after multiple retries: %w", err)
		}

		//文字列型で受け取る
		if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
			//警告コードもログ出力
			warnCode, _ := response["sWarningCode"].(string) //警告コードも文字列
			warnText, _ := response["sWarningText"].(string)
			tc.logger.Error("注文APIがエラーを返しました", zap.String("result_code", resultCode), zap.String("result_text", response["sResultText"].(string)), zap.String("warning_code", warnCode), zap.String("warning_text", warnText))
			return nil, fmt.Errorf("order API returned an error: %s - %s", resultCode, response["sResultText"])
		}
		// 注文成功時の処理 (レスポンスから必要な情報を抽出)
		order.ID = response["sOrderNumber"].(string) // 注文番号
		order.Status = "pending"                     // ステータスを更新

		return order, nil
	}
	return nil, fmt.Errorf("place order failed after multiple retries: %w", err) // 最終的なエラー
}

func (tc *TachibanaClientImple) GetOrderStatus(ctx context.Context, requestURL string, orderID string) (*domain.Order, error) {
	// リトライ処理
	var err error
	for retries := 0; retries < 3; retries++ {
		payload := map[string]string{
			"sCLMID":       clmidOrderListDetail,
			"sOrderNumber": orderID,
			"sEigyouDay":   "",                       // 必要に応じて営業日を設定
			"p_no":         tc.getPNo(),              // p_no を設定
			"p_sd_date":    formatSDDate(time.Now()), // p_sd_date を設定
		}

		response, err := sendRequest(ctx, tc, requestURL, payload)

		if err != nil {
			if isRetryableError(err) && retries < 2 {
				tc.logger.Warn("GetOrderStatus request failed, retrying...", zap.Error(err), zap.Int("retry", retries+1))
				time.Sleep(time.Duration(1+retries*2) * time.Second) // 指数バックオフ
				continue
			}
			return nil, fmt.Errorf("get order status failed after multiple retries: %w", err)
		}

		if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
			return nil, fmt.Errorf("order status API returned an error: %s - %s", resultCode, response["sResultText"])
		}

		order := &domain.Order{
			ID:     response["sOrderNumber"].(string),
			Status: response["sOrderStatus"].(string), // APIのsOrderStatusを使用
			// 他の必要なフィールドもマッピング
		}

		return order, nil
	}
	return nil, fmt.Errorf("get order status failed after multiple retries: %w", err) // 最終的なエラー
}

func (tc *TachibanaClientImple) CancelOrder(ctx context.Context, requestURL string, orderID string) error {
	// リトライ処理
	var err error

	for retries := 0; retries < 3; retries++ {
		payload := map[string]string{
			"sCLMID":          clmidCancelOrder,
			"sOrderNumber":    orderID,
			"sEigyouDay":      "", // 必要に応じて営業日を設定
			"sSecondPassword": tc.secret,
			"p_no":            tc.getPNo(),              // p_no を設定
			"p_sd_date":       formatSDDate(time.Now()), // p_sd_date を設定
		}
		response, err := sendRequest(ctx, tc, requestURL, payload)
		if err != nil {
			if isRetryableError(err) && retries < 2 {
				tc.logger.Warn("CancelOrder request failed, retrying...", zap.Error(err), zap.Int("retry", retries+1))
				time.Sleep(time.Duration(1+retries*2) * time.Second)
				continue
			}
			return fmt.Errorf("cancel order failed after multiple retries: %w", err)
		}

		if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
			return fmt.Errorf("cancel order API returned an error: %s - %s", resultCode, response["sResultText"])
		}
		return nil //成功
	}
	return fmt.Errorf("cancel order failed after multiple retries: %w", err) // 最終的なエラー
}

func isRetryableError(err error) bool {
	// ToDo: リトライ可能なエラーか判定するロジック
	return true
}
