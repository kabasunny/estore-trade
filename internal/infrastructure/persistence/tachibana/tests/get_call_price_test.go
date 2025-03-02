// internal/infrastructure/persistence/tachibana/tests/get_call_price_test.go
package tachibana_test

import (
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetCallPrice(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t) // SetupTestClient を使う

	// テストデータをセット
	callPriceMap := map[string]domain.CallPrice{
		"101": {UnitNumber: "101", Price1: 3000, UnitPrice1: 1},
		"102": {UnitNumber: "102", Price1: 10000, UnitPrice1: 5},
	}
	client.SetCallPriceMapForTest(callPriceMap)

	t.Run("Success", func(t *testing.T) {
		callPrice, ok := client.GetCallPrice("101")
		assert.True(t, ok)
		assert.Equal(t, "101", callPrice.UnitNumber)
		assert.Equal(t, 3000.0, callPrice.Price1)
		assert.Equal(t, 1.0, callPrice.UnitPrice1)
	})

	t.Run("Not Found", func(t *testing.T) {
		_, ok := client.GetCallPrice("999") // 存在しないキー
		assert.False(t, ok)
	})
}
