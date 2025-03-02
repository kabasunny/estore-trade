// internal/infrastructure/persistence/master/test_helpers.go
package master

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"strconv"
)

// GetSystemStatusFromDB は SystemStatus をDBから取得するヘルパー関数 (テスト可能)
func GetSystemStatusFromDB(ctx context.Context, db *sql.DB) (*domain.SystemStatus, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `SELECT system_status_key, login_permission, system_state FROM system_statuses` //仮
	row := db.QueryRowContext(ctx, query)

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

// GetDateInfoFromDB は DateInfo をDBから取得するヘルパー関数
func GetDateInfoFromDB(ctx context.Context, db *sql.DB) (*domain.DateInfo, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `SELECT date_key, prev_business_day1, the_day, next_business_day1, stock_delivery_date FROM date_infos` //仮
	row := db.QueryRowContext(ctx, query)

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

// GetAllCallPricesFromDB は全ての CallPrice をDBから取得するヘルパー関数
func GetAllCallPricesFromDB(ctx context.Context, db *sql.DB) ([]domain.CallPrice, error) {
	query := `SELECT unit_number, apply_date, price1, price2, price3, price4, price5, price6, price7, price8, price9, price10,
       price11, price12, price13, price14, price15, price16, price17, price18, price19, price20,
       unit_price1, unit_price2, unit_price3, unit_price4, unit_price5, unit_price6, unit_price7, unit_price8, unit_price9, unit_price10,
       unit_price11, unit_price12, unit_price13, unit_price14, unit_price15, unit_price16, unit_price17, unit_price18, unit_price19, unit_price20,
       decimal1, decimal2, decimal3, decimal4, decimal5, decimal6, decimal7, decimal8, decimal9, decimal10,
       decimal11, decimal12, decimal13, decimal14, decimal15, decimal16, decimal17, decimal18, decimal19, decimal20
		FROM call_prices` //仮

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var callPrices []domain.CallPrice
	for rows.Next() {
		var cp domain.CallPrice
		var unitNumber int
		err := rows.Scan( // ここを修正: 全てのカラムに対応する変数を指定
			&unitNumber, &cp.ApplyDate, &cp.Price1, &cp.Price2, &cp.Price3, &cp.Price4, &cp.Price5,
			&cp.Price6, &cp.Price7, &cp.Price8, &cp.Price9, &cp.Price10,
			&cp.Price11, &cp.Price12, &cp.Price13, &cp.Price14, &cp.Price15, &cp.Price16, &cp.Price17, &cp.Price18, &cp.Price19, &cp.Price20,
			&cp.UnitPrice1, &cp.UnitPrice2, &cp.UnitPrice3, &cp.UnitPrice4, &cp.UnitPrice5,
			&cp.UnitPrice6, &cp.UnitPrice7, &cp.UnitPrice8, &cp.UnitPrice9, &cp.UnitPrice10,
			&cp.UnitPrice11, &cp.UnitPrice12, &cp.UnitPrice13, &cp.UnitPrice14, &cp.UnitPrice15, &cp.UnitPrice16, &cp.UnitPrice17, &cp.UnitPrice18, &cp.UnitPrice19, &cp.UnitPrice20,
			&cp.Decimal1, &cp.Decimal2, &cp.Decimal3, &cp.Decimal4, &cp.Decimal5,
			&cp.Decimal6, &cp.Decimal7, &cp.Decimal8, &cp.Decimal9, &cp.Decimal10,
			&cp.Decimal11, &cp.Decimal12, &cp.Decimal13, &cp.Decimal14, &cp.Decimal15, &cp.Decimal16, &cp.Decimal17, &cp.Decimal18, &cp.Decimal19, &cp.Decimal20,
		)
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

// GetAllIssueMastersFromDB は全てのIssueMasterをDBから取得する(テスト用)
func GetAllIssueMastersFromDB(ctx context.Context, db *sql.DB) ([]domain.IssueMaster, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `SELECT issue_code, issue_name, trading_unit, tokutei_f FROM issue_masters` //仮
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []domain.IssueMaster
	for rows.Next() {
		var issue domain.IssueMaster
		if err := rows.Scan(&issue.IssueCode, &issue.IssueName, &issue.TradingUnit, &issue.TokuteiF); err != nil { //仮
			return nil, err
		}
		issues = append(issues, issue)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return issues, nil
}

// GetAllIssueMarketMastersFromDB はIssueMarketMasterをDBから取得する(テスト用)
func GetAllIssueMarketMastersFromDB(ctx context.Context, db *sql.DB) ([]domain.IssueMarketMaster, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `
        SELECT issue_code, market_code, price_range_min, price_range_max, sinyou_c, previous_close,
               issue_kubun_c, zyouzyou_kubun, call_price_unit_number, call_price_unit_number_yoku
        FROM issue_market_masters` // 仮のテーブル名とカラム名
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueMarkets []domain.IssueMarketMaster
	for rows.Next() {
		var im domain.IssueMarketMaster
		err := rows.Scan(
			&im.IssueCode,
			&im.MarketCode,
			&im.PriceRangeMin,
			&im.PriceRangeMax,
			&im.SinyouC,
			&im.PreviousClose,
			&im.IssueKubunC,
			&im.ZyouzyouKubun,
			&im.CallPriceUnitNumber,
			&im.CallPriceUnitNumberYoku,
		)
		if err != nil {
			return nil, err
		}
		issueMarkets = append(issueMarkets, im)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issueMarkets, nil
}

// GetAllIssueMarketRegulationsFromDB はIssueMarketRegulationをDBから取得する(テスト用)
func GetAllIssueMarketRegulationsFromDB(ctx context.Context, db *sql.DB) ([]domain.IssueMarketRegulation, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `
        SELECT issue_code, listed_market, stop_kubun, genbutu_urituke,
               seido_sinyou_sinki_uritate, ippan_sinyou_sinki_uritate, sinyou_syutyu_kubun
        FROM issue_market_regulations` // 仮のテーブル名とカラム名
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issueRegulations []domain.IssueMarketRegulation
	for rows.Next() {
		var ir domain.IssueMarketRegulation
		err := rows.Scan(
			&ir.IssueCode,
			&ir.ListedMarket,
			&ir.StopKubun,
			&ir.GenbutuUrituke,
			&ir.SeidoSinyouSinkiUritate,
			&ir.IppanSinyouSinkiUritate,
			&ir.SinyouSyutyuKubun,
		)
		if err != nil {
			return nil, err
		}
		issueRegulations = append(issueRegulations, ir)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return issueRegulations, nil
}

// GetAllOperationStatusKabuFromDB はOperationStatusKabuをDBから取得する(テスト用)
func GetAllOperationStatusKabuFromDB(ctx context.Context, db *sql.DB) ([]domain.OperationStatusKabu, error) {
	// ... (内容は前回の回答と同じ) ...
	query := `SELECT listed_market, unit, status FROM operation_statuses_kabu` // 仮のテーブル名とカラム名
	rows, err := db.QueryContext(ctx, query)                                   // QueryContext を使用
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
