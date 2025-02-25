// internal/infrastructure/persistence/order/mthd_cancel_order.go

package order

import (
	"context"
	"fmt"
	"time"
)

func (r *orderRepository) CancelOrder(ctx context.Context, orderID string) error {
	query := `
	UPDATE orders
	SET status = $2, updated_at = $3
	WHERE id = $1
`
	res, err := r.db.ExecContext(ctx, query, orderID, "canceled", time.Now()) // statusをcanceledに
	if err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
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
