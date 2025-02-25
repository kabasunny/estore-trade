package ranking

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *rankingRepository) SaveRanking(ctx context.Context, ranking []domain.Ranking) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS rankings (
    id SERIAL PRIMARY KEY,
    rank INTEGER NOT NULL,
    issue_code VARCHAR(10) NOT NULL,
    trading_value FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    `
	_, err := r.db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create rankings table: %w", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO rankings(rank, issue_code, trading_value, created_at) VALUES($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, item := range ranking {
		_, err := stmt.ExecContext(ctx, item.Rank, item.IssueCode, item.TradingValue, item.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert ranking data: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
