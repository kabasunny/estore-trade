// internal/infrastructure/persistence/tachibana/tests/get_call_price_test.go
package tachibana_test

import (
	"sync"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetCallPrice(t *testing.T) {
	// テスト用の TachibanaClientImple インスタンスを作成 (callPriceMap を初期化)
	client := &tachibana.TachibanaClientImple{
		// ... 他のフィールド ...
		CallPriceMap: map[string]domain.CallPrice{
			"101": {UnitNumber: 101, Price1: 3000, UnitPrice1: 1},
			"102": {UnitNumber: 102, Price1: 10000, UnitPrice1: 5},
		},
		Mu: sync.RWMutex{}, // テスト時に RWMutex が初期化されるようにする
	}

	t.Run("Success", func(t *testing.T) {
		callPrice, ok := client.GetCallPrice("101")
		assert.True(t, ok)
		assert.Equal(t, 101, callPrice.UnitNumber)
		assert.Equal(t, 3000.0, callPrice.Price1)
		assert.Equal(t, 1.0, callPrice.UnitPrice1)
	})

	t.Run("Not Found", func(t *testing.T) {
		_, ok := client.GetCallPrice("999") // 存在しないキー
		assert.False(t, ok)
	})
}
