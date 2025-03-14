// internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_sell_limit_test.go
package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_CreditSellLimitPayload は信用新規売り指値注文のペイロードをテストします。
func TestPlaceOrder_CreditSellLimitPayload(t *testing.T) {
	t.Run("正常系: 信用新規売り指値注文のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (信用新規売り指値注文)
		order := &domain.Order{
			Symbol:     "7974",  // 例: 任天堂
			Side:       "short", // "sell" から "short" に変更
			OrderType:  "limit",
			TradeType:  "credit_open", // 信用新規
			Quantity:   100,
			Price:      9500.0, // 指値価格
			MarketCode: "00",   // 東証
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
			"sBaibaiKubun":              "1",    // 売
			"sCondition":                "0",    // 指値
			"sOrderPrice":               "9500", // 指値価格
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "2", // 信用新規
			"sOrderExpireDay":           "0", // 当日限り
			"sGyakusasiOrderType":       "0", // 逆指値なし
			"sGyakusasiZyouken":         "0",
			"sGyakusasiPrice":           "*",
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
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_credit_sell_limit_test.go
