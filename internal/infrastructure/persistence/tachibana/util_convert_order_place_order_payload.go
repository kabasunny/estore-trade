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

	// 1. 売買区分 (sBaibaiKubun)
	var baibaiKubun string
	switch order.Side {
	case "long", "buy": // "buy" を追加
		baibaiKubun = baibaiKubunBuy
	case "short", "sell": // "sell" を追加
		baibaiKubun = baibaiKubunSell
	default:
		return nil, fmt.Errorf("invalid order side: %s", order.Side)
	}

	// 数量のバリデーション
	if order.Quantity <= 0 {
		return nil, fmt.Errorf("invalid order quantity: %d", order.Quantity)
	}

	// トリガー価格のバリデーション (逆指値注文の場合)
	if (order.OrderType == "stop" || order.OrderType == "stop_limit" ||
		order.OrderType == "credit_close_stop" || order.OrderType == "credit_close_stop_limit") &&
		order.TriggerPrice <= 0 {
		return nil, fmt.Errorf("invalid trigger price: %f", order.TriggerPrice)
	}

	// 2. 現物/信用区分 (sGenkinShinyouKubun) と、信用取引の場合は新規/返済の区別
	var genkinShinyouKubun string
	//var tatebiType string               // 建日タイプ (信用返済時に使用)  // ここはもう使わない
	var tategyokuZyoutoekiKazeiC string             // 建玉譲渡益課税区分 (信用返済時に使用)
	var aCLMKabuHensaiData []map[string]interface{} // 信用返済時の建玉指定データ

	//  OrderType から、現物/信用、新規/返済 を判定
	switch order.OrderType {
	case "market", "limit", "stop", "stop_limit": // 通常の注文（現物 or 信用新規）
		if order.TradeType == "credit_open" { // 信用新規
			genkinShinyouKubun = "2" // 信用新規
		} else { // "現物" or その他
			genkinShinyouKubun = genkinShinyouKubunGenbutsu // 現物扱い
		}
		tategyokuZyoutoekiKazeiC = "*" // デフォルト値"*"を設定

	// 信用返済
	case "credit_close_market", "credit_close_limit", "credit_close_stop", "credit_close_stop_limit":
		genkinShinyouKubun = "4" // 信用返済
		//tatebiType = "1"                  // 個別指定 (仮)  // ここももう使わない
		//tategyokuZyoutoekiKazeiC = "*" // 指定なし (仮) ここを削除

		if len(order.Positions) == 0 {
			return nil, fmt.Errorf("no positions specified for credit close order")
		}
		for _, position := range order.Positions {
			if position.Quantity <= 0 {
				return nil, fmt.Errorf("invalid position quantity: %d", position.Quantity)
			}
			if position.ID == "" {
				return nil, fmt.Errorf("invalid position ID")
			}
		}

		// 建玉情報を aCLMKabuHensaiData に変換 (order.Positions がある前提)
		for _, position := range order.Positions { // domain.OrderにPositionsフィールドが必要
			aCLMKabuHensaiData = append(aCLMKabuHensaiData, map[string]interface{}{
				"sTategyokuNumber": position.ID,                     // 建玉番号 (domain.Position から取得)
				"sTatebiZyuni":     "1",                             // 建日順 (仮)
				"sOrderSuryou":     strconv.Itoa(position.Quantity), // 数量 (domain.Position から取得)
			})
		}
		tategyokuZyoutoekiKazeiC = "*" // デフォルト値"*"を設定
	// 現引
	case "credit_to_spot":
		genkinShinyouKubun = "4" // 信用返済
		baibaiKubun = "7"        // 現引
		//tatebiType = "1"                  // 個別指定 (仮)  // ここももう使わない
		tategyokuZyoutoekiKazeiC = "1" // 指定なし (仮)
		// 建玉情報を aCLMKabuHensaiData に変換 (order.Positions がある前提)
		for _, position := range order.Positions { // domain.OrderにPositionsフィールドが必要
			aCLMKabuHensaiData = append(aCLMKabuHensaiData, map[string]interface{}{
				"sTategyokuNumber": position.ID,                     // 建玉番号 (domain.Position から取得)
				"sTatebiZyuni":     "1",                             // 建日順 (仮)
				"sOrderSuryou":     strconv.Itoa(position.Quantity), // 数量 (domain.Position から取得)
			})
		}

	// 現渡
	case "spot_to_credit":
		genkinShinyouKubun = "4"
		baibaiKubun = "5"
		//tatebiType = "1"  // ここももう使わない
		tategyokuZyoutoekiKazeiC = "1"
		for _, position := range order.Positions {
			aCLMKabuHensaiData = append(aCLMKabuHensaiData, map[string]interface{}{
				"sTategyokuNumber": position.ID,
				"sTatebiZyuni":     "1",
				"sOrderSuryou":     strconv.Itoa(position.Quantity),
			})
		}

	default:
		return nil, fmt.Errorf("unsupported order type: %s", order.OrderType)
	}

	// 3. 基本的なパラメータ
	payload := map[string]interface{}{
		"sCLMID":                    clmidPlaceOrder,
		"sZyoutoekiKazeiC":          zyoutoekiKazeiCTokutei,       // 特定口座 (固定)
		"sIssueCode":                order.Symbol,                 // 銘柄コード
		"sSizyouC":                  order.MarketCode,             //市場コード
		"sBaibaiKubun":              baibaiKubun,                  // 売買区分
		"sOrderSuryou":              strconv.Itoa(order.Quantity), // 注文数量
		"sGenkinShinyouKubun":       genkinShinyouKubun,           // 現物/信用区分
		"sOrderExpireDay":           orderExpireDay,               // 有効期限 (当日限り)
		"sSecondPassword":           tc.secret,                    // 第二パスワード
		"p_no":                      tc.getPNo(),                  // p_no
		"p_sd_date":                 formatSDDate(time.Now()),     // システム日付
		"sJsonOfmt":                 "4",                          // JSON出力フォーマット (固定)
		"sCondition":                "0",                          // 執行条件 (デフォルト: 指値)
		"sOrderPrice":               "0",                          // 注文価格 (デフォルト: 成行)  <-- ここも修正
		"sGyakusasiOrderType":       "0",                          // 逆指値注文タイプ (デフォルト: なし)
		"sGyakusasiZyouken":         "0",                          // 逆指値条件 (デフォルト: なし)
		"sGyakusasiPrice":           "*",                          // 逆指値価格 (デフォルト: *)
		"sTatebiType":               "*",                          //  "*"（現物または新規）
		"sTategyokuZyoutoekiKazeiC": tategyokuZyoutoekiKazeiC,     // 建玉譲渡益課税区分 (信用返済時のみ)
	}

	// ★★★ 以下の部分を 3. 基本的なパラメータ の直後に移動 ★★★
	if genkinShinyouKubun == "4" {
		payload["sTatebiType"] = "1" // 信用返済
		if baibaiKubun == "7" {      //現引
			payload["sOrderPrice"] = "*"
		}
	}
	// ★★★ 移動ここまで ★★★

	// 4. 信用返済注文の場合は、建玉指定データを追加
	if len(aCLMKabuHensaiData) > 0 {
		payload["aCLMKabuHensaiData"] = aCLMKabuHensaiData
	}

	// 5. 注文タイプに応じたパラメータ設定 (成行、指値、逆指値)
	switch order.OrderType {
	case "market", "credit_close_market", "credit_to_spot", "spot_to_credit":
		// 成行注文: sCondition, sOrderPrice, sGyakusasiOrderType, sGyakusasiZyouken, sGyakusasiPrice はデフォルト値のまま

	case "limit", "credit_close_limit":
		// 指値注文
		payload["sCondition"] = conditionSashine                               // 指値
		payload["sOrderPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64) // 指値価格

	case "stop", "credit_close_stop":
		// 通常逆指値
		payload["sCondition"] = conditionSashine                                            // 指値
		payload["sOrderPrice"] = "*"                                                        // 注文価格は "*"
		payload["sGyakusasiOrderType"] = "1"                                                // 通常逆指値
		payload["sGyakusasiZyouken"] = strconv.FormatFloat(order.TriggerPrice, 'f', -1, 64) // 逆指値条件
		//以下修正
		if order.Price == 0 {
			payload["sGyakusasiPrice"] = "0"
		} else {
			payload["sGyakusasiPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64) // 逆指値価格
		}

	case "stop_limit", "credit_close_stop_limit":
		// 通常+逆指値
		payload["sCondition"] = conditionSashine                                            // 指値
		payload["sOrderPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64)              // 通常注文の価格
		payload["sGyakusasiOrderType"] = "2"                                                // 通常+逆指値
		payload["sGyakusasiZyouken"] = strconv.FormatFloat(order.TriggerPrice, 'f', -1, 64) // 逆指値条件

		//payload["sGyakusasiPrice"] = strconv.FormatFloat(order.Price, 'f', -1, 64)          // 逆指値価格　ここが間違っていた
		if order.AfterTriggerOrderType == "market" {
			payload["sGyakusasiPrice"] = "0"
		} else if order.AfterTriggerOrderType == "limit" {
			payload["sGyakusasiPrice"] = strconv.FormatFloat(order.AfterTriggerPrice, 'f', -1, 64)
		} else {
			// エラー処理 (AfterTriggerOrderType が不正な値の場合)
			return nil, fmt.Errorf("invalid AfterTriggerOrderType: %s", order.AfterTriggerOrderType)
		}
	default: // 現状、ここには来ないはずだが、念のためエラー処理
		return nil, fmt.Errorf("unsupported order type: %s", order.OrderType)
	}

	return payload, nil
}
