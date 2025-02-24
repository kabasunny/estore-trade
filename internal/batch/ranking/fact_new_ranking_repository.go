package ranking

import (
	"database/sql"
)

func NewRankingRepository(db *sql.DB) *rankingRepository { // interfaceではなく、構造体を返す
	return &rankingRepository{db: db}
}
