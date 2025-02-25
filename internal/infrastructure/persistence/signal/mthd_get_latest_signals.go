// internal/infrastructure/persistence/signal/mthd_get_latest_signals.go
package signal

import (
	"context"
	"estore-trade/internal/domain"
)

func (r *signalRepository) GetLatestSignals(ctx context.Context, limit int) ([]domain.Signal, error) {
	query := `
        SELECT id, issue_code, side, priority, created_at
        FROM signals
        ORDER BY created_at DESC
        LIMIT $1
    `
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var signals []domain.Signal
	for rows.Next() {
		var signal domain.Signal
		if err := rows.Scan(&signal.ID, &signal.IssueCode, &signal.Side, &signal.Priority, &signal.CreatedAt); err != nil {
			return nil, err
		}
		signals = append(signals, signal)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return signals, nil
}
