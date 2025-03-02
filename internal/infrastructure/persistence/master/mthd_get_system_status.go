// internal/infrastructure/persistence/master/mthd_get_system_status.go
package master

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

// SystemStatusをDBから取得
func (r *masterDataRepository) getSystemStatus(ctx context.Context) (*domain.SystemStatus, error) {
	query := `SELECT system_status_key, login_permission, system_state FROM system_statuses` //仮
	row := r.db.QueryRowContext(ctx, query)

	var systemStatus domain.SystemStatus
	err := row.Scan(&systemStatus.SystemStatusKey, &systemStatus.LoginPermission, &systemStatus.SystemState)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // レコードがない場合はnilを返す
		}
		return nil, err
	}
	return &systemStatus, nil
}
