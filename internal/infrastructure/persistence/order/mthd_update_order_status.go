package order

import (
	"context"
	"fmt"
	"time"
)

// 指定された注文IDの注文ステータスを更新するメソッド
func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	// クエリを定義：ordersテーブルの指定されたIDの注文のステータスと更新日時を更新
	query := `
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `
	// クエリを実行し、指定された値をプレースホルダーに挿入
	res, err := r.db.ExecContext(ctx, query, orderID, status, time.Now())
	if err != nil {
		// クエリ実行中にエラーが発生した場合、そのエラーを返す
		return fmt.Errorf("failed to update order status: %w", err)
	}

	// 影響を受けた行数を取得
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		// 行数の取得中にエラーが発生した場合、そのエラーを返す
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	// 影響を受けた行数が0の場合、注文が見つからなかったことを示すエラーを返す
	if rowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}

	// 正常に更新された場合、nilを返す
	return nil
}
