// internal/infrastructure/persistence/tachibana/tests/get_issue_master_test.go
package tachibana_test

import (
	"context"
	"sync" // Mutexの初期化に必要
	"testing"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetIssueMaster(t *testing.T) {
	// テスト用の TachibanaClientImple インスタンスを作成
	client := &tachibana.TachibanaClientImple{
		// 各フィールドを初期化
		IssueMap: map[string]domain.IssueMaster{
			"7974": {IssueCode: "7974", IssueName: "任天堂"},
			"9984": {IssueCode: "9984", IssueName: "ソフトバンク"},
		},
		// 他の必要なフィールドも初期化
		CallPriceMap:             make(map[string]domain.CallPrice),
		IssueMarketMap:           make(map[string]map[string]domain.IssueMarketMaster),
		IssueMarketRegulationMap: make(map[string]map[string]domain.IssueMarketRegulation),
		OperationStatusKabuMap:   make(map[string]map[string]domain.OperationStatusKabu),
		Mu:                       sync.RWMutex{}, // RWMutex の初期化
		TargetIssueCodesMu:       sync.RWMutex{}, // RWMutex の初期化
	}

	t.Run("Success", func(t *testing.T) {
		issue, ok := client.GetIssueMaster("7974")
		assert.True(t, ok)
		assert.Equal(t, "7974", issue.IssueCode)
		assert.Equal(t, "任天堂", issue.IssueName)
	})

	t.Run("Not Found", func(t *testing.T) {
		_, ok := client.GetIssueMaster("9999") // 存在しないキー
		assert.False(t, ok)
	})

	t.Run("Target Issue - Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定
		client.Mu.Lock() // グローバルなミューテックスをロック
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"7974"})
		client.TargetIssueCodesMu.Unlock()
		client.Mu.Unlock() // グローバルなミューテックスをアンロック

		issue, ok := client.GetIssueMaster("7974")
		assert.True(t, ok)
		assert.Equal(t, "7974", issue.IssueCode)
	})

	t.Run("Target Issue - Not Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定 (7974 は含まない)
		client.Mu.Lock()
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"9984"})
		client.TargetIssueCodesMu.Unlock()
		client.Mu.Unlock()
		_, ok := client.GetIssueMaster("7974") // ターゲット銘柄リストに含まれていない
		assert.False(t, ok)
	})

	t.Run("Target Issue - Empty", func(t *testing.T) {
		client.Mu.Lock()
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{}) //空にする
		client.TargetIssueCodesMu.Unlock()
		client.Mu.Unlock()
		_, ok := client.GetIssueMaster("7974")
		assert.False(t, ok) //ターゲット銘柄がない場合は、false
	})
}
