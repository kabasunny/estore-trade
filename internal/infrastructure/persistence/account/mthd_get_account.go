// internal/infrastructure/persistence/account/mthd_get_account.go
package account

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

func (r *accountRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	query := `
        SELECT id, user_id, account_type, balance, available_balance, margin, created_at, updated_at
        FROM accounts
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)

	account := &domain.Account{}
	err := row.Scan(
		&account.ID,
		&account.UserID,
		&account.AccountType,
		&account.Balance,
		&account.AvailableBalance,
		&account.Margin,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	positions, err := r.getPositions(ctx, id)
	if err != nil {
		return nil, err
	}
	account.Positions = positions

	return account, nil
}
