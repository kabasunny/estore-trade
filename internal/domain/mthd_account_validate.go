// internal/domain/strct_account.go
package domain

import (
	"fmt"
)

// Validate Account 構造体のバリデーション
func (a *Account) Validate() error {
	if a.ID == "" {
		return fmt.Errorf("account ID is required")
	}
	if a.UserID == "" {
		return fmt.Errorf("user ID is required")
	}
	if a.AccountType != "special" {
		return fmt.Errorf("invalid account type: %s", a.AccountType)
	}
	if a.Balance < 0 {
		return fmt.Errorf("balance cannot be negative")
	}
	// 他のバリデーションルール...
	return nil
}
