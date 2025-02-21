package account

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

// 指定されたIDのアカウントをデータベースから取得するメソッド
func (r *accountRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	// クエリを定義：指定されたIDのアカウントを取得
	query := `
        SELECT id, balance, created_at, updated_at
        FROM accounts
        WHERE id = $1
    `
	// クエリを実行し、結果の行を取得
	row := r.db.QueryRowContext(ctx, query, id) // $1はプレースホルダ で　idに対応

	// 結果を格納するための Account インスタンスを作成
	account := &domain.Account{}

	// 行からデータをスキャンして Account インスタンスに格納
	err := row.Scan(&account.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		// 指定されたIDのアカウントが見つからなかった場合
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err // その他のエラーが発生した場合、エラーを返す
	}

	// ポジションの取得
	positions, err := r.getPositions(ctx, id)
	if err != nil {
		return nil, err // ポジションの取得中にエラーが発生した場合、エラーを返す
	}
	account.Positions = positions // 取得したポジションをアカウントに格納

	return account, nil // アカウント情報を返す
}
