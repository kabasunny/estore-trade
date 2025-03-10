// internal/infrastructure/persistence/tachibana/util_convert_order_place_order_payload.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
	"strconv"
	"time"
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

	// 基本的なパラメータ (すべての注文タイプで共通) + 成行注文のデフォルト値
	payload := map[string]interface{}{
		"sCLMID":                    clmidPlaceOrder,
		"sZyoutoekiKazeiC":          zyoutoekiKazeiCTokutei,
		"sIssueCode":                order.Symbol,
		"sSizyouC":                  sizyouCToushou,
		"sBaibaiKubun":              baibaiKubun,
		"sOrderSuryou":              strconv.Itoa(order.Quantity),
		"sGenkinShinyouKubun":       genkinShinyouKubunGenbutsu, // 現物
		"sOrderExpireDay":           orderExpireDay,
		"sSecondPassword":           tc.secret,
		"p_no":                      tc.getPNo(),
		"p_sd_date":                 formatSDDate(time.Now()),
		"sJsonOfmt":                 "4",
		"sCondition":                "0", // 成行注文のデフォルト値
		"sOrderPrice":               "0", // 成行注文のデフォルト値
		"sGyakusasiOrderType":       "0", // 成行注文のデフォルト値
		"sGyakusasiZyouken":         "0", // 成行注文のデフォルト値
		"sGyakusasiPrice":           "*", // 成行注文のデフォルト値
		"sTatebiType":               "*", // 成行注文のデフォルト値
		"sTategyokuZyoutoekiKazeiC": "*", // 成行注文のデフォルト値
	}

	// 執行条件と価格 (注文の種類によって異なる)
	fmt.Printf("\n order.OrderType: %s\n\n", order.OrderType) // デバッグ
	switch order.OrderType {
	case "market":
		// 成行注文 (デフォルト値が設定済みなので、ここでは特に何もしない)
		// payload["sCondition"] = "0"  // 不要 (デフォルト値)
		// payload["sOrderPrice"] = "0" // 不要 (デフォルト値)

	case "limit":
		// 指値注文
		payload["sCondition"] = conditionSashine
		payload["sOrderPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64)

	case "stop":
		// 通常逆指値
		payload["sCondition"] = conditionSashine
		payload["sOrderPrice"] = "*" // sOrderPriceは*
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

	// デバッグ出力: payload の内容を出力
	fmt.Printf("DEBUG: payload = %v\n", payload)

	return payload, nil
}
