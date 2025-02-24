// internal/infrastructure/persistence/signal/fact_new_signal_repository.go
package signal

import (
	"database/sql"
	"estore-trade/internal/domain"
)

func NewSignalRepository(db *sql.DB) domain.SignalRepository { // SignalRepository インターフェースを返す
	return &signalRepository{db: db}
}
