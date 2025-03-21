package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_SpotSellStopLimitPayload は現物売り逆指値注文(通常+逆指値)のペイロードをテストします。
func TestPlaceOrder_SpotSellStopLimitPayload(t *testing.T) {
	t.Run("正常系: 現物売り逆指値注文 (通常+逆指値) のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (現物売り逆指値注文)
		order := &domain.Order{
			Symbol:                "7974", // 例: 任天堂
			Side:                  "short",
			OrderType:             "stop_limit",
			TradeType:             "", // 現物
			Quantity:              100,
			Price:                 8000.0,  //逆指値価格
			TriggerPrice:          9500.0,  // 逆指値トリガー価格
			MarketCode:            "00",    // 東証
			AfterTriggerOrderType: "limit", // 追加
			AfterTriggerPrice:     8000.0,  // 追加
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
			"sBaibaiKubun":              "1",    // 売り
			"sCondition":                "0",    // 指値
			"sOrderPrice":               "8000", //逆指値価格
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "0",    // 現物
			"sOrderExpireDay":           "0",    // 当日限り
			"sGyakusasiOrderType":       "2",    // 通常+逆指値
			"sGyakusasiZyouken":         "9500", // 逆指値トリガー価格
			"sGyakusasiPrice":           "8000", // 逆指値価格
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

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_spot_sell_stop_limit_test.go
