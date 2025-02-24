// internal/infrastructure/persistence/ranking/ranking_repository.go
package ranking

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"
)

// ランキングデータをDBに保存する
func (r *rankingRepository) SaveRanking(ctx context.Context, ranking []domain.Ranking) error {
	// 1. まず、テーブルの存在を確認し、なければ作成する
	// (マイグレーションツールなどを使っている場合は不要)
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
	// 2. データを挿入
	tx, err := r.db.BeginTx(ctx, nil) // トランザクション開始
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // ロールバック (Commit されなかった場合)

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
	if err := tx.Commit(); err != nil { // トランザクションをコミット
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

// 必要に応じて、GetRanking, GetLatestRanking などのメソッドを実装
