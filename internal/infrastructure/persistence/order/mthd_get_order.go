package order

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

// 指定されたIDの注文をデータベースから取得するメソッド
func (r *orderRepository) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	// クエリを定義：ordersテーブルから指定されたIDの注文を取得
	query := `
        SELECT id, symbol, order_type, side, quantity, price, status, tachibana_order_id, created_at, updated_at
        FROM orders
        WHERE id = $1
    `
	// クエリを実行し、結果の行を取得
	row := r.db.QueryRowContext(ctx, query, id)

	// 結果を格納するための Order インスタンスを作成
	order := &domain.Order{}

	// 行からデータをスキャンして Order インスタンスに格納
	err := row.Scan(
		&order.ID, &order.Symbol, &order.OrderType, &order.Side,
		&order.Quantity, &order.Price, &order.Status, &order.TachibanaOrderID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		// 指定されたIDの注文が見つからなかった場合
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err // その他のエラーが発生した場合、エラーを返す
	}
	return order, nil // 注文情報を返す
}
