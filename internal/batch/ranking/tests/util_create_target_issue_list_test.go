package ranking_test

import (
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTargetIssueList(t *testing.T) {
	// テストデータの準備
	rankingData := []domain.Ranking{
		{Rank: 1, IssueCode: "1001"},
		{Rank: 2, IssueCode: "1002"},
		{Rank: 3, IssueCode: "1003"},
		{Rank: 4, IssueCode: "1004"},
	}

	t.Run("Limit within ranking size", func(t *testing.T) {
		targetIssues := ranking.CreateTargetIssueList(rankingData, 2) // 上位2件

		// 取得した件数が正しいことを確認
		assert.Len(t, targetIssues, 2)

		// 正しい銘柄コードが取得できていることを確認
		assert.Equal(t, "1001", targetIssues[0].IssueCode)
		assert.Equal(t, "1002", targetIssues[1].IssueCode)
	})

	t.Run("Limit exceeds ranking size", func(t *testing.T) {
		targetIssues := ranking.CreateTargetIssueList(rankingData, 10) // ランキングのサイズを超える件数を指定

		// ランキングのサイズと同じ件数が取得されることを確認
		assert.Len(t, targetIssues, len(rankingData))

		// 正しい銘柄コードが取得できていることを確認
		assert.Equal(t, "1001", targetIssues[0].IssueCode)
		assert.Equal(t, "1002", targetIssues[1].IssueCode)
		assert.Equal(t, "1003", targetIssues[2].IssueCode)
		assert.Equal(t, "1004", targetIssues[3].IssueCode)
	})

	t.Run("Limit is 0", func(t *testing.T) {
		targetIssues := ranking.CreateTargetIssueList(rankingData, 0) // 0件

		// 空のスライスが返されることを確認
		assert.Empty(t, targetIssues)
	})
}
