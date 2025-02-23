// internal/handler/trading.go
package handler

import (
	"fmt"

	"estore-trade/internal/domain"
)

func validateOrderRequest(order *domain.Order) error {
	if order.Quantity <= 0 {
		return fmt.Errorf("invalid order quantity: %d", order.Quantity)
	}
	// 他のバリデーションルール...
	return nil
}
