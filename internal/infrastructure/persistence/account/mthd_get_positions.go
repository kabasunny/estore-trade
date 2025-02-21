package account

import (
	"context"
	"estore-trade/internal/domain"
)

// 指定されたアカウントIDに対応するポジションをデータベースから取得するメソッド
func (r *accountRepository) getPositions(ctx context.Context, accountID string) ([]domain.Position, error) {
	// クエリを定義：指定されたアカウントIDに対応するポジションを取得
	query := `
        SELECT symbol, quantity, price
        FROM positions
        WHERE account_id = $1
    `
	// クエリを実行し、結果の行を取得
	rows, err := r.db.QueryContext(ctx, query, accountID) // $1はプレースホルダ で　accountIDに対応
	if err != nil {
		return nil, err // エラーが発生した場合、エラーを返す
	}
	defer rows.Close() // 処理が終わった後に結果セットを閉じる

	// ポジションのスライスを初期化
	var positions []domain.Position
	for rows.Next() {
		var position domain.Position
		// 行からデータをスキャンしてポジションに格納
		if err := rows.Scan(&position.Symbol, &position.Quantity, &position.Price); err != nil {
			return nil, err // スキャン中にエラーが発生した場合、エラーを返す
		}
		positions = append(positions, position) // ポジションをスライスに追加
	}
	// ループの最後にエラーが発生したかを確認
	if err := rows.Err(); err != nil {
		return nil, err // エラーが発生した場合、エラーを返す
	}
	return positions, nil // 取得したポジションのスライスを返す
}
