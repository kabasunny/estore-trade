// internal/infrastructure/persistence/master/mthd_get_all_operation_status_kabu.go
package master

import (
	"context"
	"estore-trade/internal/domain"
)

// OperationStatusKabuをDBから取得
func (r *masterDataRepository) getAllOperationStatusKabu(ctx context.Context) ([]domain.OperationStatusKabu, error) {
	query := `SELECT listed_market, unit, status FROM operation_statuses_kabu` // 仮のテーブル名とカラム名
	rows, err := r.db.QueryContext(ctx, query)                                 // QueryContext を使用
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operationStatuses []domain.OperationStatusKabu
	for rows.Next() {
		var os domain.OperationStatusKabu
		if err := rows.Scan(&os.ListedMarket, &os.Unit, &os.Status); err != nil {
			return nil, err
		}
		operationStatuses = append(operationStatuses, os)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return operationStatuses, nil
}
