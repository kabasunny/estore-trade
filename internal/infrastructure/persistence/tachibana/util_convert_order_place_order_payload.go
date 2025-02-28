// internal/infrastructure/persistence/tachibana/util_convert_order_place_order_payload.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
	"strconv"
)

// ConvertOrderToPlaceOrderPayload は、domain.Order を立花証券の注文リクエストペイロードに変換
func ConvertOrderToPlaceOrderPayload(order *domain.Order, tc *TachibanaClientImple) (map[string]interface{}, error) {
	var baibaiKubun string
	switch order.Side {
	case "buy":
		baibaiKubun = baibaiKubunBuy
	case "sell":
		baibaiKubun = baibaiKubunSell
	default:
		return nil, fmt.Errorf("invalid order side: %s", order.Side)
	}
	payload := map[string]interface{}{
		"sCLMID":              clmidPlaceOrder,
		"sZyoutoekiKazeiC":    zyoutoekiKazeiCTokutei,
		"sIssueCode":          order.Symbol,
		"sSizyouC":            sizyouCToushou,
		"sBaibaiKubun":        baibaiKubun, // ここに baibaiKubun を設定
		"sOrderSuryou":        strconv.Itoa(order.Quantity),
		"sGenkinShinyouKubun": genkinShinyouKubunGenbutsu,
		"sOrderExpireDay":     orderExpireDay,
		"sSecondPassword":     tc.Secret,
		"p_no":                tc.getPNo(),
		//"p_sd_date":           formatSDDate(time.Now()), // formatSDDate は util_format_sd_date.go で定義
	}
	// 執行条件と価格
	switch order.OrderType {
	case "market":
		payload["sCondition"] = conditionSashine // 0: 指定なし (立花証券の仕様)
		payload["sOrderPrice"] = "0"             // 成行
	case "limit":
		payload["sCondition"] = conditionSashine
		payload["sOrderPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64)
	case "stop":
		// 通常逆指値
		payload["sCondition"] = conditionSashine
		payload["sOrderPrice"] = "*"
		payload["sGyakusasiOrderType"] = "1"
		payload["sGyakusasiZyouken"] = strconv.FormatFloat(order.TriggerPrice, 'f', -1, 64)
		payload["sGyakusasiPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64)
	case "stop_limit":
		// 通常+逆指値
		payload["sCondition"] = conditionSashine
		payload["sOrderPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64) //通常注文の価格
		payload["sGyakusasiOrderType"] = "2"
		payload["sGyakusasiZyouken"] = strconv.FormatFloat(order.TriggerPrice, 'f', -1, 64)
		payload["sGyakusasiPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64) //逆指値になった場合の価格

	default:
		return nil, fmt.Errorf("unsupported order type: %s", order.OrderType)
	}

	return payload, nil
}
