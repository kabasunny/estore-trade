package ranking

import (
	"database/sql"
	"estore-trade/internal/domain" // 追加
)

func NewRankingRepository(db *sql.DB) domain.RankingRepository { // 戻り値を interface に変更
	return &rankingRepository{db: db}
}
