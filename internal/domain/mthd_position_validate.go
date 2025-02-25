// internal/domain/strct_position.go
package domain

import (
	"fmt"
)

// Validate Position 構造体のバリデーション
func (p *Position) Validate() error {
	if p.Symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if p.Side != "long" && p.Side != "short" {
		return fmt.Errorf("invalid side: %s", p.Side)
	}
	if p.Quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if p.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	return nil
}
