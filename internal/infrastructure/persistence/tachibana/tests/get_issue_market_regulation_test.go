// internal/infrastructure/persistence/tachibana/tests/get_issue_market_regulation_test.go
package tachibana_test

import (
	"context"
	"testing"

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetIssueMarketRegulation(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)
	defer httpmock.DeactivateAndReset()

	t.Logf("Before Test: %+v", client.IssueMarketRegulationMap)

	t.Run("Success", func(t *testing.T) {
		regulation, ok := client.GetIssueMarketRegulation("7974", "00")
		t.Logf("Regulation: %+v", regulation)
		assert.True(t, ok)
		assert.Equal(t, "7974", regulation.IssueCode)
		assert.Equal(t, "00", regulation.ListedMarket)
	})

	t.Run("Issue Not Found", func(t *testing.T) {
		_, ok := client.GetIssueMarketRegulation("9999", "00") // 存在しない銘柄コード
		assert.False(t, ok)
	})

	t.Run("Market Not Found", func(t *testing.T) {
		_, ok := client.GetIssueMarketRegulation("7974", "XX") // 存在しない市場コード
		assert.False(t, ok)
	})

	t.Run("Target Issue - Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"7974"})
		client.TargetIssueCodesMu.Unlock()

		regulation, ok := client.GetIssueMarketRegulation("7974", "00")
		assert.True(t, ok)
		assert.Equal(t, "7974", regulation.IssueCode)
		assert.Equal(t, "00", regulation.ListedMarket)
	})

	t.Run("Target Issue - Not Found", func(t *testing.T) {
		// ターゲット銘柄リストを設定 (7974 は含まない)
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{"9984"})
		client.TargetIssueCodesMu.Unlock()
		_, ok := client.GetIssueMarketRegulation("7974", "00") // ターゲット銘柄リストに含まれていない
		assert.False(t, ok)
	})

	t.Run("Target Issue - Empty", func(t *testing.T) {
		client.TargetIssueCodesMu.Lock()
		client.SetTargetIssues(context.Background(), []string{}) //空にする
		client.TargetIssueCodesMu.Unlock()
		_, ok := client.GetIssueMarketRegulation("7974", "00")
		assert.False(t, ok) //ターゲット銘柄がない場合は、false
	})
}
