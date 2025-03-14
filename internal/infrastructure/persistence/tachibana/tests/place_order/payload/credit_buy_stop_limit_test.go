// internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_buy_stop_limit_test.go
package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_CreditBuyStopLimitPayload は信用新規買い通常+逆指値注文のペイロードをテストします。
func TestPlaceOrder_CreditBuyStopLimitPayload(t *testing.T) {
	t.Run("正常系: 信用新規買い通常+逆指値注文(トリガー後指値)のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (信用新規買い通常+逆指値注文)
		order := &domain.Order{
			Symbol:                "7974",        // 例: 任天堂
			Side:                  "long",        // "buy" から "long" に変更
			OrderType:             "stop_limit",  // 通常+逆指値
			TradeType:             "credit_open", // 信用新規
			Quantity:              100,
			Price:                 10000.0, // 逆指値価格（指値）
			TriggerPrice:          9500.0,  // 逆指値トリガー価格
			MarketCode:            "00",    // 東証
			AfterTriggerOrderType: "limit", // 追加
			AfterTriggerPrice:     9800.0,  // 追加: トリガー後の指値
		}

		// ConvertOrderToPlaceOrderPayload 関数を呼び出してペイロードを生成
		payload, err := tachibana.ConvertOrderToPlaceOrderPayload(order, tc)
		assert.NoError(t, err)

		// 期待されるペイロード
		expectedPayload := map[string]interface{}{
			"sCLMID":                    "CLMKabuNewOrder",
			"sZyoutoekiKazeiC":          "1", // 特定口座
			"sIssueCode":                "7974",
			"sSizyouC":                  "00",
			"sBaibaiKubun":              "3",     // 買
			"sCondition":                "0",     // 指値
			"sOrderPrice":               "10000", // 逆指値価格（通常注文の価格）
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "2",    // 信用新規
			"sOrderExpireDay":           "0",    // 当日限り
			"sGyakusasiOrderType":       "2",    // 通常+逆指値
			"sGyakusasiZyouken":         "9500", // 逆指値トリガー価格
			"sGyakusasiPrice":           "9800", // 逆指値価格: これを AfterTriggerPrice に合わせる
			"sTatebiType":               "*",
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(),
			"p_no":                      tc.GetPNoForTest(),
			"p_sd_date":                 payload["p_sd_date"], // 動的なのでpayloadから取得
			"sJsonOfmt":                 "4",
		}

		// 生成されたペイロードを検証
		assert.Equal(t, expectedPayload, payload)
	})

	t.Run("正常系: 信用新規買い通常+逆指値注文(トリガー後成行)のペイロードが正しく生成されること", func(t *testing.T) {
		tc := tachibana.CreateTestClient(t, nil)
		order := &domain.Order{
			Symbol:                "7974",
			Side:                  "long",
			OrderType:             "stop_limit",
			TradeType:             "credit_open",
			Quantity:              100,
			Price:                 9000.0, //トリガー前の指値
			TriggerPrice:          9500.0,
			MarketCode:            "00",
			AfterTriggerOrderType: "market", // 追加: トリガー後成行
		}

		payload, err := tachibana.ConvertOrderToPlaceOrderPayload(order, tc)
		assert.NoError(t, err)

		expectedPayload := map[string]interface{}{
			"sCLMID":                    "CLMKabuNewOrder",
			"sZyoutoekiKazeiC":          "1",
			"sIssueCode":                "7974",
			"sSizyouC":                  "00",
			"sBaibaiKubun":              "3",
			"sCondition":                "0",    // 指値
			"sOrderPrice":               "9000", //トリガー前の指値
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "2",
			"sOrderExpireDay":           "0",
			"sGyakusasiOrderType":       "2", // 通常+逆指値
			"sGyakusasiZyouken":         "9500",
			"sGyakusasiPrice":           "0", //トリガー後、成行
			"sTatebiType":               "*",
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(),
			"p_no":                      tc.GetPNoForTest(),
			"p_sd_date":                 payload["p_sd_date"],
			"sJsonOfmt":                 "4",
		}
		assert.Equal(t, expectedPayload, payload)
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_buy_stop_limit_test.go
