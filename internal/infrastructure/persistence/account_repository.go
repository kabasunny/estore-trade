// internal/infrastructure/persistence/account_repository.go
package persistence

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"time"
)

type accountRepository struct {
	db *sql.DB
}

// NewAccountRepository は AccountRepository の新しいインスタンスを作成する
func NewAccountRepository(db *sql.DB) domain.AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) GetAccount(ctx context.Context, id string) (*domain.Account, error) {
	query := `
        SELECT id, balance, created_at, updated_at
        FROM accounts
        WHERE id = $1
    `
	row := r.db.QueryRowContext(ctx, query, id)
	account := &domain.Account{}
	err := row.Scan(&account.ID, &account.Balance, &account.CreatedAt, &account.UpdatedAt) // &account.UserID,
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Account not found.
		}
		return nil, err
	}

	// ポジションの取得 (別途関数に切り出す方が良い)
	positions, err := r.getPositions(ctx, id)
	if err != nil {
		return nil, err
	}
	account.Positions = positions

	return account, nil
}

func (r *accountRepository) UpdateAccount(ctx context.Context, account *domain.Account) error {
	query := `
        UPDATE accounts
        SET balance = $2, updated_at = $3
        WHERE id = $1
    `
	_, err := r.db.ExecContext(ctx, query, account.ID, account.Balance, time.Now()) // &account.UserID,
	return err
}

// getPositions はアカウントのポジションを取得 (accountRepository のメソッド)
func (r *accountRepository) getPositions(ctx context.Context, accountID string) ([]domain.Position, error) {
	query := `
        SELECT symbol, quantity, price
        FROM positions
        WHERE account_id = $1
    `
	rows, err := r.db.QueryContext(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []domain.Position
	for rows.Next() {
		var position domain.Position
		if err := rows.Scan(&position.Symbol, &position.Quantity, &position.Price); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return positions, nil
}
