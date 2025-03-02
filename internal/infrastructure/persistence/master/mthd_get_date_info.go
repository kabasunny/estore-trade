// internal/infrastructure/persistence/master/mthd_get_date_info.go
package master

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
)

// DateInfoをDBから取得
func (r *masterDataRepository) getDateInfo(ctx context.Context) (*domain.DateInfo, error) {
	query := `SELECT date_key, prev_business_day1, the_day, next_business_day1, stock_delivery_date FROM date_infos` //仮
	row := r.db.QueryRowContext(ctx, query)

	var dateInfo domain.DateInfo
	err := row.Scan(&dateInfo.DateKey, &dateInfo.PrevBusinessDay1, &dateInfo.TheDay, &dateInfo.NextBusinessDay1, &dateInfo.StockDeliveryDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &dateInfo, nil
}
