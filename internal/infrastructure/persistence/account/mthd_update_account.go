package account

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

// 指定されたアカウントのデータを更新するメソッド
func (r *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	// クエリを定義：指定されたアカウントIDの残高と更新日時を更新
	query := `
        UPDATE accounts
        SET balance = $2, updated_at = $3
        WHERE id = $1
    `
	// クエリを実行し、アカウントの残高と更新日時を更新
	_, err := r.db.ExecContext(ctx, query, account.ID, account.Balance, time.Now()) // $1,$2,$3はプレースホルダ で　account.ID, account.Balance, time.Now()に対応
	// クエリ実行中にエラーが発生した場合、そのエラーを返す
	return err
}
