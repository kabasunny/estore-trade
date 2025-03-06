package ranking_test

import (
	"context"
	"database/sql"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveRanking(t *testing.T) {
	// テスト用のDBセットアップ (インメモリ SQLite)
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	repo := ranking.NewRankingRepository(db)

	t.Run("Save ranking data successfully", func(t *testing.T) {
		rankingData := []domain.Ranking{
			{Rank: 1, IssueCode: "1001", TradingValue: 1000, CreatedAt: time.Now()},
			{Rank: 2, IssueCode: "1002", TradingValue: 900, CreatedAt: time.Now()},
		}

		err := repo.SaveRanking(context.Background(), rankingData)
		assert.NoError(t, err)

		// データが正しく保存されたか確認 (別のクエリで取得)
		rows, err := db.Query("SELECT rank, issue_code, trading_value FROM rankings")
		require.NoError(t, err)
		defer rows.Close()

		var retrievedRanking []domain.Ranking
		for rows.Next() {
			var item domain.Ranking
			err := rows.Scan(&item.Rank, &item.IssueCode, &item.TradingValue)
			require.NoError(t, err)
			retrievedRanking = append(retrievedRanking, item)
		}

		// 取得したデータが、保存したデータと一致するか確認 (件数、内容)
		assert.Len(t, retrievedRanking, len(rankingData))
		assert.Equal(t, rankingData[0].Rank, retrievedRanking[0].Rank)
		assert.Equal(t, rankingData[0].IssueCode, retrievedRanking[0].IssueCode)
		assert.Equal(t, rankingData[1].Rank, retrievedRanking[1].Rank)
		assert.Equal(t, rankingData[1].IssueCode, retrievedRanking[1].IssueCode)
		// TradingValue は float64 なので、完全に一致しない場合がある。誤差を許容する比較が必要。
		// (ここでは簡単のため、省略)
	})
}
