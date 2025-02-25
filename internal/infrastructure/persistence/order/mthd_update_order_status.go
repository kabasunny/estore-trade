// internal/infrastructure/persistence/order/mthd_update_order_status.go
package order

import (
	"context"
	"fmt"
	"time"
)

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status string) error {
	query := `
        UPDATE orders
        SET status = $2, updated_at = $3
        WHERE id = $1
    `
	res, err := r.db.ExecContext(ctx, query, orderID, status, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("order not found: %s", orderID)
	}

	return nil
}
