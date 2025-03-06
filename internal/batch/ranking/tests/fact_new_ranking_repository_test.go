package ranking_test

import (
	"database/sql"
	"estore-trade/internal/batch/ranking"
	"estore-trade/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRankingRepository(t *testing.T) {
	// ダミーのDB接続を作成 (実際にはテスト用のDBを用意する)
	db, err := sql.Open("sqlite3", ":memory:") // インメモリの SQLite DB を使用
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	repo := ranking.NewRankingRepository(db)

	// 戻り値が期待される型であること、nil でないことを確認
	assert.NotNil(t, repo)
	assert.Implements(t, (*domain.RankingRepository)(nil), repo) // domain.RankingRepository を実装しているか
}
