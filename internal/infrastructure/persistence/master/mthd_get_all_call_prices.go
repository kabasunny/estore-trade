// internal/infrastructure/persistence/master/mthd_get_all_call_prices.go
package master

import (
	"context"
	"estore-trade/internal/domain"
	"strconv"
)

// CallPriceをDBから取得
func (r *masterDataRepository) getAllCallPrices(ctx context.Context) ([]domain.CallPrice, error) {
	query := `SELECT unit_number, apply_date, price1, price2, price3, price4, price5, price6, price7, price8, price9, price10,
       price11, price12, price13, price14, price15, price16, price17, price18, price19, price20,
       unit_price1, unit_price2, unit_price3, unit_price4, unit_price5, unit_price6, unit_price7, unit_price8, unit_price9, unit_price10,
       unit_price11, unit_price12, unit_price13, unit_price14, unit_price15, unit_price16, unit_price17, unit_price18, unit_price19, unit_price20,
       decimal1, decimal2, decimal3, decimal4, decimal5, decimal6, decimal7, decimal8, decimal9, decimal10,
       decimal11, decimal12, decimal13, decimal14, decimal15, decimal16, decimal17, decimal18, decimal19, decimal20
		FROM call_prices` //仮

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var callPrices []domain.CallPrice
	for rows.Next() {
		var cp domain.CallPrice
		var unitNumber int
		err := rows.Scan(
			&unitNumber, &cp.ApplyDate, &cp.Price1, &cp.Price2, &cp.Price3, &cp.Price4, &cp.Price5,
			&cp.Price6, &cp.Price7, &cp.Price8, &cp.Price9, &cp.Price10,
			&cp.Price11, &cp.Price12, &cp.Price13, &cp.Price14, &cp.Price15, &cp.Price16, &cp.Price17, &cp.Price18, &cp.Price19, &cp.Price20,
			&cp.UnitPrice1, &cp.UnitPrice2, &cp.UnitPrice3, &cp.UnitPrice4, &cp.UnitPrice5,
			&cp.UnitPrice6, &cp.UnitPrice7, &cp.UnitPrice8, &cp.UnitPrice9, &cp.UnitPrice10,
			&cp.UnitPrice11, &cp.UnitPrice12, &cp.UnitPrice13, &cp.UnitPrice14, &cp.UnitPrice15, &cp.UnitPrice16, &cp.UnitPrice17, &cp.UnitPrice18, &cp.UnitPrice19, &cp.UnitPrice20,
			&cp.Decimal1, &cp.Decimal2, &cp.Decimal3, &cp.Decimal4, &cp.Decimal5,
			&cp.Decimal6, &cp.Decimal7, &cp.Decimal8, &cp.Decimal9, &cp.Decimal10,
			&cp.Decimal11, &cp.Decimal12, &cp.Decimal13, &cp.Decimal14, &cp.Decimal15, &cp.Decimal16, &cp.Decimal17, &cp.Decimal18, &cp.Decimal19, &cp.Decimal20,
		) //仮
		if err != nil {
			return nil, err
		}
		cp.UnitNumber = strconv.Itoa(unitNumber) //intからstringに
		callPrices = append(callPrices, cp)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return callPrices, nil
}
