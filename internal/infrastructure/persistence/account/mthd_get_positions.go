// internal/infrastructure/persistence/account/mthd_get_positions.go
package account

import (
	"context"
	"estore-trade/internal/domain"
)

func (r *accountRepository) getPositions(ctx context.Context, accountID string) ([]domain.Position, error) {
	query := `
        SELECT symbol, quantity, price, side
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
		if err := rows.Scan(&position.Symbol, &position.Quantity, &position.Price, &position.Side); err != nil {
			return nil, err
		}
		positions = append(positions, position)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return positions, nil
}
