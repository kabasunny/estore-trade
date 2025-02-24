package signal

import (
	"database/sql"
)

type signalRepository struct {
	db *sql.DB
}
