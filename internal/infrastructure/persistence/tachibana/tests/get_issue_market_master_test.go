// internal/infrastructure/persistence/tachibana/tests/get_issue_market_master_test.go
package tachibana_test

import (
	"context"
	"sync"
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetIssueMarketMaster(t *testing.T) {
	// テスト用の TachibanaClientImple インスタンスを作成 (issueMarketMap を初期化)
	client := &tachibana.TachibanaClientImple{
		IssueMarketMap: map[string]map[string]domain.IssueMarketMaster{
			"7974": {
				"00": {IssueCode: "7974", MarketCode: "00"}, // 東証
			},
			"9984": {
				"00": {IssueCode: "9984", MarketCode: "00"}, // 東証
			},
		},
		// 他の必要なフィールドも初期化
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMap:                 make(map[string]domain.IssueMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
		Mu:                       sync.RWMutex{}, // RWMutex の初期化
		TargetIssueCodesMu:       sync.RWMutex{},
	}

	t.Run("Success", func(t *testing.T) {
		issueMarket, ok := client.GetIssueMarketMaster("7974", "00")
		assert.True(t, ok)
		assert.Equal(t, "7974", issueMarket.IssueCode)
		assert.Equal(t, "00", issueMarket.MarketCode)
	})

	t.Run("Issue Not Found", func(t *testing.T) {
		_, ok := client.GetIssueMarketMaster("9999", "00") // 存在しない銘柄コード
		assert.False(t, ok)
	})

	t.Run("Market Not Found", func(t *testing.T) {
		_, ok := client.GetIssueMarketMaster("7974", "XX") // 存在しない市場コード
		assert.False(t, ok)
	})
	t.Run("Target Issue - Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"7974"})
		client.TargetIssueCodesMu.Unlock()

		issueMarket, ok := client.GetIssueMarketMaster("7974", "00")
		assert.True(t, ok)
		assert.Equal(t, "7974", issueMarket.IssueCode)
		assert.Equal(t, "00", issueMarket.MarketCode)
	})

	t.Run("Target Issue - Not Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定 (7974 は含まない)
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"9984"})
		client.TargetIssueCodesMu.Unlock()
		_, ok := client.GetIssueMarketMaster("7974", "00") // ターゲット銘柄リストに含まれていない
		assert.False(t, ok)
	})

	t.Run("Target Issue - Empty", func(t *testing.T) {
		client.Mu.Lock()
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{}) //空にする
		client.TargetIssueCodesMu.Unlock()
		client.Mu.Unlock()
		issueMarket, ok := client.GetIssueMarketMaster("7974", "00") //issueMarketMapにデータが存在する
		assert.True(t, ok)                                           //ターゲット銘柄がない場合は、true
		assert.Equal(t, "7974", issueMarket.IssueCode)
		assert.Equal(t, "00", issueMarket.MarketCode)
	})
}
