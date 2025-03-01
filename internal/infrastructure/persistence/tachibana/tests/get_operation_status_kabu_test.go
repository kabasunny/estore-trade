// internal/infrastructure/persistence/tachibana/tests/get_operation_status_kabu_test.go
package tachibana_test

import (
	//"sync" //不要
	"testing"

	//"estore-trade/internal/domain"//不要
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetOperationStatusKabu(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)

	t.Run("Success", func(t *testing.T) {
		status, ok := client.GetOperationStatusKabu("00", "0101")
		assert.True(t, ok)
		assert.Equal(t, "00", status.ListedMarket)
		assert.Equal(t, "0101", status.Unit)
		assert.Equal(t, "001", status.Status)
	})

	t.Run("Market Not Found", func(t *testing.T) {
		_, ok := client.GetOperationStatusKabu("XX", "0101") // 存在しない市場
		assert.False(t, ok)
	})

	t.Run("Unit Not Found", func(t *testing.T) {
		_, ok := client.GetOperationStatusKabu("00", "9999") // 存在しない単位
		assert.False(t, ok)
	})
}
