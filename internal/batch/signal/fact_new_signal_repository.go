// internal/infrastructure/persistence/signal/signal_repository.go
package signal

import (
	"database/sql"
)

func NewSignalRepository(db *sql.DB) *signalRepository {
	return &signalRepository{db: db}
}
