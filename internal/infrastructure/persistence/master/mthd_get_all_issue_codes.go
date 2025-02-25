// internal/infrastructure/persistence/master/mthd_get_all_issue_codes.go
package master

import (
	"context"
)

func (r *masterDataRepository) GetAllIssueCodes(ctx context.Context) ([]string, error) {
	query := `
        SELECT issue_code
        FROM issue_masters
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueCodes []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		issueCodes = append(issueCodes, code)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return issueCodes, nil
}
