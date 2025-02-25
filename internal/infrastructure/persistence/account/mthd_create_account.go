// internal/infrastructure/persistence/account/mthd_create_account.go
package account

import (
	"context"
	"estore-trade/internal/domain"
	"time"
)

func (r *accountRepository) CreateAccount(ctx context.Context, account *domain.Account) error {
	query := `
        INSERT INTO accounts (id, user_id, account_type, balance, available_balance, margin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err := r.db.ExecContext(ctx, query,
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		time.Now(),
		time.Now(),
	)
	return err
}
