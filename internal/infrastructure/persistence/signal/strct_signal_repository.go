// internal/infrastructure/persistence/signal/strct_signal_repository.go
package signal

import (
	"database/sql"
)

type signalRepository struct {
	db *sql.DB
}
