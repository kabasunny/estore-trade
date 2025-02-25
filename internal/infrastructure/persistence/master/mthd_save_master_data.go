// internal/infrastructure/persistence/master/mthd_save_master_data.go
package master

import (
	"context"
	"estore-trade/internal/domain" // 追加
	"fmt"
)

func (r *masterDataRepository) SaveMasterData(ctx context.Context, m *domain.MasterData) error {
	// 1. トランザクション開始
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // エラー発生時はロールバック

	// 2. テーブルの存在確認と作成 (マイグレーションツールを使用する場合は不要)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS issue_masters (
        issue_code VARCHAR(10) PRIMARY KEY,
        issue_name VARCHAR(255) NOT NULL,
        trading_unit INTEGER NOT NULL,
        tokutei_f BOOLEAN NOT NULL
    );
    `
	_, err = tx.ExecContext(ctx, createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create issue_masters table: %w", err)
	}

	// 3. 既存データの削除 (必要に応じて。今回は毎回全データを入れ替える想定)
	_, err = tx.ExecContext(ctx, "DELETE FROM issue_masters")
	if err != nil {
		return fmt.Errorf("failed to delete existing issue_masters data: %w", err)
	}

	// 4. データの挿入 (issue_master のみ)
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO issue_masters(issue_code, issue_name, trading_unit, tokutei_f) VALUES($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, issue := range m.GetIssueMap() { // 変更: GetIssueMap() を使用
		var tokuteiF int
		if issue.TokuteiF == "1" {
			tokuteiF = 1
		} else {
			tokuteiF = 0
		}
		_, err := stmt.ExecContext(ctx, issue.IssueCode, issue.IssueName, issue.TradingUnit, tokuteiF)
		if err != nil {
			return fmt.Errorf("failed to insert issue_master data: %w", err)
		}
	}

	// 5. トランザクションのコミット
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
