// internal/infrastructure/persistence/tachibana/tests/get_issue_master_test.go
package tachibana_test

import (
	//"context"
	//"sync" // Mutexの初期化に必要 //不要
	"testing"

	//"estore-trade/internal/domain" //不要
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetIssueMaster(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)

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

}

// TestTargetIssuesGetIssueMaster にまとめる
func TestTargetIssuesGetIssueMaster(t *testing.T) {
	client, _ := tachibana.SetupTestClient(t)

	tests := []struct {
		name       string
		issueCodes []string
		issueCode  string
		expect     bool
	}{
		{
			name:       "Target Issue - Found",
			issueCodes: []string{"7974"},
			issueCode:  "7974",
			expect:     true,
		},
		{
			name:       "Target Issue - Not Found",
			issueCodes: []string{"9984"},
			issueCode:  "7974",
			expect:     false,
		},
		{
			name:       "Target Issue - Empty",
			issueCodes: []string{},
			issueCode:  "7974",
			expect:     false, // ターゲットリストが空の場合は、issueMapにあってもfalse
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.SetTargetIssueCodesForTest(tt.issueCodes) //ヘルパーメソッドを使用
			_, ok := client.GetIssueMaster(tt.issueCode)
			assert.Equal(t, tt.expect, ok)
		})
	}
}
