// internal/infrastructure/persistence/master/mthd_get_issue_master.go
package master

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"fmt"
)

func (r *masterDataRepository) GetIssueMaster(ctx context.Context, issueCode string) (*domain.IssueMaster, error) {
	query := `
		SELECT issue_code, issue_name, trading_unit, tokutei_f
		FROM issue_masters
		WHERE issue_code = $1
	`
	row := r.db.QueryRowContext(ctx, query, issueCode)

	var issueMaster domain.IssueMaster // ポインタをやめる
	var tokuteiF int                   // Assuming TokuteiF is stored as an integer (0 or 1)
	err := row.Scan(
		&issueMaster.IssueCode,
		&issueMaster.IssueName,
		&issueMaster.TradingUnit,
		&tokuteiF, // Scan into the integer
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	// Convert the integer back to a string representation for consistency
	issueMaster.TokuteiF = fmt.Sprintf("%d", tokuteiF)

	return &issueMaster, nil
}
