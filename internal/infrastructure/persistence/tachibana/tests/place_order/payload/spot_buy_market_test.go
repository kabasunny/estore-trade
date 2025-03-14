package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_SpotBuyMarketPayload は現物買い成行注文のペイロードをテストします。
func TestPlaceOrder_SpotBuyMarketPayload(t *testing.T) {
	t.Run("正常系: 現物買い成行注文のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (現物買い成行注文)
		order := &domain.Order{
			Symbol:     "7974", // 例: 任天堂
			Side:       "long", // "buy" から "long" に変更
			OrderType:  "market",
			TradeType:  "", // 現物
			Quantity:   100,
			MarketCode: "00", // 東証
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
			"sBaibaiKubun":              "3", // 買
			"sCondition":                "0", // 指定なし (成行)
			"sOrderPrice":               "0", // 成行
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "0", // 現物
			"sOrderExpireDay":           "0", // 当日限り
			"sGyakusasiOrderType":       "0", // 逆指値なし
			"sGyakusasiZyouken":         "0",
			"sGyakusasiPrice":           "*",
			"sTatebiType":               "*",
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(), // GetPasswordForTest() を使用
			"p_no":                      tc.GetPNoForTest(),      // GetPNoForTest() を使用
			"p_sd_date":                 payload["p_sd_date"],    //動的な値なので、生成されたペイロードから取得
			"sJsonOfmt":                 "4",
		}

		// 生成されたペイロードを検証
		assert.Equal(t, expectedPayload, payload)
	})
}

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_spot_buy_market_test.go
