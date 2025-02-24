package ranking

import (
	"database/sql"
)

type rankingRepository struct {
	db *sql.DB
}
