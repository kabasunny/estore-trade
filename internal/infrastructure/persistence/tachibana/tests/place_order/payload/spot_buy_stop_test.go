package payload_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

// TestPlaceOrder_SpotBuyStopPayload は現物買い逆指値注文のペイロードをテストします。
func TestPlaceOrder_SpotBuyStopPayload(t *testing.T) {
	t.Run("正常系: 現物買い逆指値注文 (通常) のペイロードが正しく生成されること", func(t *testing.T) {
		// テスト用の TachibanaClientImple インスタンスを作成
		tc := tachibana.CreateTestClient(t, nil)

		// domain.Order オブジェクトを作成 (現物買い逆指値注文)
		order := &domain.Order{
			Symbol:       "7974", // 例: 任天堂
			Side:         "long",
			OrderType:    "stop",
			TradeType:    "", // 現物
			Quantity:     100,
			Price:        10000.0, // 逆指値価格（指値）
			TriggerPrice: 9500.0,  // 逆指値トリガー価格
			MarketCode:   "00",    // 東証
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
			"sCondition":                "0", // 指値
			"sOrderPrice":               "*", //逆指値の時は"*"
			"sOrderSuryou":              "100",
			"sGenkinShinyouKubun":       "0",     // 現物
			"sOrderExpireDay":           "0",     // 当日限り
			"sGyakusasiOrderType":       "1",     // 通常逆指値
			"sGyakusasiZyouken":         "9500",  // 逆指値トリガー価格
			"sGyakusasiPrice":           "10000", //逆指値価格
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

// go test -v ./internal/infrastructure/persistence/tachibana/tests/payload/place_order_spot_buy_stop_test.go
