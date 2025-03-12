// internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_close_buy_stop_limit_test.go
package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_CreditCloseBuyStopLimitPayload は信用返済買い通常+逆指値注文のペイロードをテストします。
func TestPlaceOrder_CreditCloseBuyStopLimitPayload(t *testing.T) {
	t.Run("正常系: 信用返済買い通常+逆指値注文(トリガー後指値)のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (信用返済買い通常+逆指値注文)
		order := &domain.Order{
			Symbol:                "7974",                    // 例: 任天堂
			Side:                  "long",                    // 買い
			OrderType:             "credit_close_stop_limit", // 信用返済 (通常+逆指値)
			Quantity:              100,
			Price:                 10000.0, // 通常注文の価格（指値）
			TriggerPrice:          9500.0,  // 逆指値トリガー価格
			MarketCode:            "00",    // 東証
			AfterTriggerOrderType: "limit",
			AfterTriggerPrice:     9800.0, //トリガー後の指値
			Positions: []domain.Position{ // 返済する建玉の情報
				{
					ID:       "202007220000402", // 建玉番号
					Quantity: 100,               // 返済数量
				},
			},
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
			"sOrderPrice":               "10000", // 通常注文の価格
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "4",    // 信用返済
			"sOrderExpireDay":           "0",    // 当日限り
			"sGyakusasiOrderType":       "2",    // 通常+逆指値
			"sGyakusasiZyouken":         "9500", // 逆指値トリガー価格
			"sGyakusasiPrice":           "9800", // 逆指値価格
			"sTatebiType":               "1",    // 個別指定
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(),
			"p_no":                      tc.GetPNoForTest(),
			"p_sd_date":                 payload["p_sd_date"], // 動的なのでpayloadから取得
			"sJsonOfmt":                 "4",
			"aCLMKabuHensaiData": []map[string]interface{}{ // 返済建玉リスト
				{
					"sTategyokuNumber": "202007220000402", // 建玉番号
					"sTatebiZyuni":     "1",               // 建日順位 (仮)
					"sOrderSuryou":     "100",             // 注文数量
				},
			},
		}

		// 生成されたペイロードを検証
		assert.Equal(t, expectedPayload, payload)
	})

	t.Run("正常系: 信用返済買い通常+逆指値注文(トリガー後成行)のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (信用返済買い通常+逆指値注文 - トリガー後成行)
		order := &domain.Order{
			Symbol:                "7974",                    // 例: 任天堂
			Side:                  "long",                    // 買い
			OrderType:             "credit_close_stop_limit", // 信用返済 (通常+逆指値)
			Quantity:              100,
			Price:                 9000.0,   // 通常注文の価格
			TriggerPrice:          9500.0,   // 逆指値トリガー価格
			MarketCode:            "00",     // 東証
			AfterTriggerOrderType: "market", // トリガー後成行
			Positions: []domain.Position{ // 返済する建玉の情報
				{
					ID:       "202007220000402", // 建玉番号
					Quantity: 100,               // 返済数量
				},
			},
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
			"sBaibaiKubun":              "3",    // 買
			"sCondition":                "0",    // 指値
			"sOrderPrice":               "9000", // 通常注文の価格
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "4",    // 信用返済
			"sOrderExpireDay":           "0",    // 当日限り
			"sGyakusasiOrderType":       "2",    // 通常+逆指値
			"sGyakusasiZyouken":         "9500", // 逆指値トリガー価格
			"sGyakusasiPrice":           "0",    // トリガー後成行
			"sTatebiType":               "1",    // 個別指定
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(),
			"p_no":                      tc.GetPNoForTest(),
			"p_sd_date":                 payload["p_sd_date"], // 動的なのでpayloadから取得
			"sJsonOfmt":                 "4",
			"aCLMKabuHensaiData": []map[string]interface{}{ // 返済建玉リスト
				{
					"sTategyokuNumber": "202007220000402", // 建玉番号
					"sTatebiZyuni":     "1",               // 建日順位 (仮)
					"sOrderSuryou":     "100",             // 注文数量
				},
			},
		}

		// 生成されたペイロードを検証
		assert.Equal(t, expectedPayload, payload)
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_close_buy_stop_limit_test.go
