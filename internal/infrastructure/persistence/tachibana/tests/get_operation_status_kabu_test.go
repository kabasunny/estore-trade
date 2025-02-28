// internal/infrastructure/persistence/tachibana/tests/get_operation_status_kabu_test.go
package tachibana_test

import (
	"sync"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetOperationStatusKabu(t *testing.T) {
	// テスト用の TachibanaClientImple インスタンスを作成 (OperationStatusKabuMap を初期化)
	client := &tachibana.TachibanaClientImple{
		OperationStatusKabuMap: map[string]map[string]domain.OperationStatusKabu{
			"00": { // 東証
				"0101": {ListedMarket: "00", Unit: "0101", Status: "001"}, // 株式
			},
		},
		// 他の必要なフィールドも初期化
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		Mu:                       sync.RWMutex{}, // RWMutex の初期化
		TargetIssueCodesMu:       sync.RWMutex{},
	}

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
