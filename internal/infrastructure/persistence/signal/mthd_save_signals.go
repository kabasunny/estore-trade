// internal/infrastructure/persistence/signal/mthd_save_signals.go
package signal

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

// SaveSignals はシグナルデータをDBに保存する
func (r *signalRepository) SaveSignals(ctx context.Context, signals []domain.Signal) error {
	// 1. テーブルの存在を確認し、なければ作成 (マイグレーションツールを使う場合は不要)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS signals (
        id SERIAL PRIMARY KEY,
        issue_code VARCHAR(10) NOT NULL,
        side VARCHAR(4) NOT NULL,
        priority INTEGER NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    `
	_, err := r.db.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create signals table: %w", err)
	}
	// 2. トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 3. データを挿入
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO signals(issue_code, side, priority, created_at) VALUES($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, signal := range signals {
		_, err := stmt.ExecContext(ctx, signal.IssueCode, signal.Side, signal.Priority, signal.CreatedAt)
		if err != nil {
			// ロールバックは defer で行われるので、ここではエラーを返すだけで良い
			return fmt.Errorf("failed to insert signal data: %w", err)
		}
	}
	// 4. トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
