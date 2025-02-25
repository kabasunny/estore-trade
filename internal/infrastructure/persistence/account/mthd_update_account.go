// internal/infrastructure/persistence/account/mthd_update_account.go
package account

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

func (r *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	query := `
        UPDATE accounts
        SET user_id = $2, account_type = $3, balance = $4, available_balance = $5, margin = $6, updated_at = $7
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		time.Now(),
	)
	return err
}
