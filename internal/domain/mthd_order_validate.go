// internal/domain/strct_order.go
package domain

import (
	"fmt"
)

// Validate Order 構造体のバリデーション
func (o *Order) Validate() error {
	if o.UUID == "" {
		return fmt.Errorf("order ID is required")
	}
	if o.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if o.Side != "long" && o.Side != "short" {
		return fmt.Errorf("invalid side: %s", o.Side)
	}
	if o.OrderType != "market" && o.OrderType != "limit" && o.OrderType != "stop" && o.OrderType != "stop_limit" {
		return fmt.Errorf("invalid order type: %s", o.OrderType)
	}
	if o.OrderType == "limit" || o.OrderType == "stop_limit" {
		if o.Price <= 0 {
			return fmt.Errorf("price must be positive for limit/stop limit orders")
		}
	}
	if o.OrderType == "stop" || o.OrderType == "stop_limit" {
		if o.TriggerPrice <= 0 {
			return fmt.Errorf("trigger price must be positive for stop/stop limit orders")
		}
	}
	if o.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	// 他のバリデーションルール...
	return nil
}
