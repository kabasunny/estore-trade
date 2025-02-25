// internal/infrastructure/persistence/signal/fact_new_signal_repository.go
package signal

import (
	"database/sql"
	"estore-trade/internal/domain"
)

func NewSignalRepository(db *sql.DB) domain.SignalRepository {
	return &signalRepository{db: db}
}
