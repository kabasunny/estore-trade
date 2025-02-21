package order

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

// 指定された注文のデータをデータベースで更新するメソッド
func (r *orderRepository) UpdateOrder(ctx context.Context, order *domain.Order) error {
	// クエリを定義：ordersテーブルの指定されたIDの注文のデータを更新
	query := `
        UPDATE orders
        SET symbol = $2, order_type = $3, side = $4, quantity = $5, price = $6, status = $7, tachibana_order_id = $8, updated_at = $9
        WHERE id = $1
    `
	// クエリを実行し、指定された値をプレースホルダーに挿入
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.Symbol, order.OrderType, order.Side,
		order.Quantity, order.Price, order.Status, order.TachibanaOrderID, time.Now(),
	)
	// クエリ実行中にエラーが発生した場合、そのエラーを返す
	return err
}
