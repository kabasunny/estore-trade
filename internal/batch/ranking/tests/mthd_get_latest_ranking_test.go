package ranking_test

import (
	"context"
	"database/sql"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite3 ドライバをインポート
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLatestRanking(t *testing.T) {
	// テスト用のDBセットアップ (インメモリ SQLite)
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// リポジトリの作成
	repo := ranking.NewRankingRepository(db)

	// テーブルの作成 (SaveRanking でテーブルが作成されるが、念のため)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS rankings (
        id SERIAL PRIMARY KEY,
        rank INTEGER NOT NULL,
        issue_code VARCHAR(10) NOT NULL,
        trading_value FLOAT NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    `
	_, err = db.Exec(createTableSQL)
	require.NoError(t, err)

	// テストデータの準備
	now := time.Now()
	rankingData := []domain.Ranking{
		{Rank: 1, IssueCode: "1001", TradingValue: 1000, CreatedAt: now.Add(-time.Hour)},  // 古いデータ
		{Rank: 2, IssueCode: "1002", TradingValue: 900, CreatedAt: now.Add(-time.Hour)},   // 古いデータ
		{Rank: 3, IssueCode: "2001", TradingValue: 2000, CreatedAt: now},                  // 新しいデータ
		{Rank: 2, IssueCode: "2002", TradingValue: 1900, CreatedAt: now},                  // 新しいデータ
		{Rank: 1, IssueCode: "2003", TradingValue: 1800, CreatedAt: now},                  // 新しいデータ
		{Rank: 4, IssueCode: "2004", TradingValue: 1700, CreatedAt: now.Add(time.Minute)}, // 最新のデータ
	}
	// テストデータ投入
	err = repo.SaveRanking(context.Background(), rankingData)
	require.NoError(t, err)

	// テスト実行
	t.Run("Get latest ranking", func(t *testing.T) {
		latestRanking, err := repo.GetLatestRanking(context.Background(), 3) // 上位3件を取得
		require.NoError(t, err)

		// 取得した件数が正しいか
		assert.Len(t, latestRanking, 3)

		// 最新のランキングが取得できているか (created_at と rank の両方で確認)

		// 1番目は created_at が最も新しい "2004"
		assert.Equal(t, "2004", latestRanking[0].IssueCode)
		assert.Equal(t, 4, latestRanking[0].Rank)

		// 2番目と3番目は created_at が同じ (now)。rank で順序が決まる。
		assert.Equal(t, "2003", latestRanking[1].IssueCode)
		assert.Equal(t, 1, latestRanking[1].Rank)

		assert.Equal(t, "2002", latestRanking[2].IssueCode)
		assert.Equal(t, 2, latestRanking[2].Rank)
	})

	t.Run("Get latest ranking with limit 0", func(t *testing.T) {
		latestRanking, err := repo.GetLatestRanking(context.Background(), 0) // limitを0
		require.NoError(t, err)
		assert.Len(t, latestRanking, 0) //何も取得されない
	})

}
