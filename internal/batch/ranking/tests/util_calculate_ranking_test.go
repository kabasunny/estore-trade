package ranking_test

import (
	"context"
	"estore-trade/internal/batch/ranking"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRanking(t *testing.T) {
	// テストデータの準備 (モックデータ)
	marketData := []ranking.MarketDataItem{ // 修正: marketDataItem を使用
		{IssueCode: "1001", Price: 100, Volume: 10}, // 売買代金: 1000
		{IssueCode: "1002", Price: 200, Volume: 5},  // 売買代金: 1000
		{IssueCode: "1003", Price: 50, Volume: 30},  // 売買代金: 1500
		{IssueCode: "1004", Price: 300, Volume: 2},  // 売買代金: 600
	}

	// コンテキストの準備 (必要に応じて)
	ctx := context.Background()

	// CalculateRanking 関数の呼び出し
	rankingData, err := ranking.CalculateRanking(ctx, marketData) //clientを削除

	// エラーがないことを確認
	assert.NoError(t, err)

	// ランキングの件数が正しいことを確認
	assert.Len(t, rankingData, 4)

	// ランキングが売買代金の降順に並んでいることを確認
	assert.Equal(t, 1, rankingData[0].Rank)
	assert.Equal(t, "1003", rankingData[0].IssueCode) // 売買代金: 1500

	assert.Equal(t, 2, rankingData[1].Rank)                                // 1000で同値
	assert.Contains(t, []string{"1001", "1002"}, rankingData[1].IssueCode) // 売買代金: 1000 (同値は順不同)

	assert.Equal(t, 3, rankingData[2].Rank)                                // 1000で同値
	assert.Contains(t, []string{"1001", "1002"}, rankingData[2].IssueCode) // 売買代金: 1000 (同値は順不同)

	assert.Equal(t, 4, rankingData[3].Rank)
	assert.Equal(t, "1004", rankingData[3].IssueCode) // 売買代金: 600

	assert.IsType(t, time.Time{}, rankingData[0].CreatedAt)
}
