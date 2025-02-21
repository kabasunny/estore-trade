package order

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

// 新しい注文をデータベースに挿入するメソッド
func (r *orderRepository) CreateOrder(ctx context.Context, order *domain.Order) error {
	// クエリを定義：ordersテーブルに新しい注文を挿入
	query := `
        INSERT INTO orders (id, symbol, order_type, side, quantity, price, status, tachibana_order_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	// クエリを実行し、指定された値をプレースホルダーに挿入
	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.Symbol, order.OrderType, order.Side,
		order.Quantity, order.Price, order.Status, order.TachibanaOrderID, time.Now(), time.Now(),
	)
	// クエリ実行中にエラーが発生した場合、そのエラーを返す
	return err
}
