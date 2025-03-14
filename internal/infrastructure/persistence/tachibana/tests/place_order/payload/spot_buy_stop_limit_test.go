// internal/infrastructure/persistence/tachibana/tests/payload/place_order_spot_buy_stop_limit_test.go
package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_SpotBuyStopLimitPayload は現物買い逆指値注文(通常+逆指値)のペイロードをテストします。
func TestPlaceOrder_SpotBuyStopLimitPayload(t *testing.T) {
	t.Run("正常系: 現物買い逆指値注文 (通常+逆指値) のペイロードが正しく生成されること", func(t *testing.T) {
		//テストスト用のTachibanaClientImpleインスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		//domain.Orderオブジェクトを作成（現物買い逆指値注文）
		order := &domain.Order{
			Symbol:                "7974",
			Side:                  "long",
			OrderType:             "stop_limit", // 通常＋逆指値
			TradeType:             "",
			Quantity:              100,
			Price:                 10000.0, //逆指値価格
			TriggerPrice:          9500.0,  //逆指値トリガー価格
			MarketCode:            "00",
			AfterTriggerOrderType: "limit", // 追加
			AfterTriggerPrice:     9800.0,  // 追加
		}
		// ConvertOrderToPlaceOrderPayload 関数を呼び出してペイロードを生成
		payload, err := tachibana.ConvertOrderToPlaceOrderPayload(order, tc)
		assert.NoError(t, err)

		//期待されるペイロード
		expectedPayload := map[string]interface{}{
			"sCLMID":                    "CLMKabuNewOrder",
			"sZyoutoekiKazeiC":          "1",
			"sIssueCode":                "7974",
			"sSizyouC":                  "00",
			"sBaibaiKubun":              "3",
			"sCondition":                "0",
			"sOrderPrice":               "10000", //逆指値価格
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "0",
			"sOrderExpireDay":           "0",
			"sGyakusasiOrderType":       "2", //通常+逆指値
			"sGyakusasiZyouken":         "9500",
			"sGyakusasiPrice":           "9800", // 逆指値価格　ここを修正
			"sTatebiType":               "*",
			"sTategyokuZyoutoekiKazeiC": "*",
			"sSecondPassword":           tc.GetPasswordForTest(),
			"p_no":                      tc.GetPNoForTest(),
			"p_sd_date":                 payload["p_sd_date"], // 動的なのでpayloadから取得
			"sJsonOfmt":                 "4",
		}
		//生成されたペイロードを検証
		assert.Equal(t, expectedPayload, payload)
	})
}
