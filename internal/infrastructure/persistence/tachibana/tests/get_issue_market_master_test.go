// internal/infrastructure/persistence/tachibana/tests/get_issue_market_master_test.go
package tachibana_test

import (
	//"context"
	//"sync"
	"testing"

	//"estore-trade/internal/domain" //いらない
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetIssueMarketMaster(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)

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
}

func TestTargetIssuesGetIssueMarketMaster(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)

	tests := []struct {
		name       string
		issueCodes []string
		issueCode  string
		marketCode string
		expect     bool
	}{
		{
			name:       "Target Issue - Found",
			issueCodes: []string{"7974"},
			issueCode:  "7974",
			marketCode: "00",
			expect:     true,
		},
		{
			name:       "Target Issue - Not Found",
			issueCodes: []string{"9984"},
			issueCode:  "7974",
			marketCode: "00",
			expect:     false,
		},
		{
			name:       "Target Issue - Empty",
			issueCodes: []string{},
			issueCode:  "7974",
			marketCode: "00",
			expect:     true, // ターゲットリストが空の場合は、issueMarketMapにあればtrue
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.SetTargetIssueCodesForTest(tt.issueCodes) //ヘルパーメソッドを使用
			_, ok := client.GetIssueMarketMaster(tt.issueCode, tt.marketCode)
			assert.Equal(t, tt.expect, ok)
		})
	}
}
