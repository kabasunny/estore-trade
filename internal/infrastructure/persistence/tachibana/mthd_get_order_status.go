// internal/infrastructure/persistence/tachibana/mthd_get_order_status.go
package tachibana

import (
	"context"
	"encoding/json"
	"estore-trade/internal/domain"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// GetOrderStatus は注文のステータスを取得
func (tc *TachibanaClientImple) GetOrderStatus(ctx context.Context, orderID string, orderDate string) (*domain.Order, error) { //orderDate string
	payload := map[string]string{
		"sCLMID":       clmidOrderListDetail,
		"sOrderNumber": orderID,
		"sEigyouDay":   orderDate, // 変数名を変更
		"p_no":         tc.getPNo(),
		"p_sd_date":    formatSDDate(time.Now()),
		"sJsonOfmt":    "4", // JSON出力フォーマット (追加)
	}

	// ペイロードをJSONにエンコード (Logout メソッドに合わせる)
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// URL エンコード
	encodedPayload := url.QueryEscape(string(payloadJSON))
	requestURL := tc.requestURL + "?" + encodedPayload

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil) // GET
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json") // なくても良い
	req, cancel := withContextAndTimeout(req, 60*time.Second)
	defer cancel()

	response, err := sendRequest(req, 3) // 3回リトライ設定し送信
	if err != nil {
		return nil, fmt.Errorf("get order status failed: %w", err)
	}

	if resultCode, ok := response["sResultCode"].(string); ok && resultCode != "0" {
		return nil, fmt.Errorf("order status API returned an error: %s - %s", resultCode, response["sResultText"])
	}

	// レスポンスから必要な情報を抽出して、domain.Order 構造体にマッピング
	order := &domain.Order{
		TachibanaOrderID: orderID, //引数のorderIDをセット
	}

	if status, ok := response["sOrderStatus"].(string); ok { //注文ステータス
		order.Status = status
	}

	if baibaiKubun, ok := response["sOrderBaibaiKubun"].(string); ok { //売買区分
		switch baibaiKubun {
		case "1":
			order.Side = "short" //ドメイン層の定数に変更
		case "3":
			order.Side = "long" //ドメイン層の定数に変更
		default:
			return nil, fmt.Errorf("invalid sOrderBaibaiKubun in order status: %s", baibaiKubun)
		}
	}

	// 約定数量 (sYakujouSuryou) を取得
	if quantityStr, ok := response["sYakujouSuryou"].(string); ok {
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sYakujouSuryou: %w", err)
		}
		order.Quantity = quantity
	}
	//sOrderOrderSuryou	注文株数,sOrderCurrentSuryou	有効株数は、必要になったら追加

	// 必要であれば、約定単価 (sYakujouPrice) なども取得
	if priceStr, ok := response["sYakujouPrice"].(string); ok {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sYakujouPrice: %w", err)
		}
		order.Price = price
	}

	return order, nil
}
