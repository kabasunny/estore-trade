// internal/infrastructure/persistence/master/tests/get_master_data_test.go
package master_test

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/master"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMasterDataRepository_GetMasterData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db)

	tests := []struct {
		name        string
		mockSetup   func(mock sqlmock.Sqlmock)
		expected    *domain.MasterData
		expectedErr error
	}{
		{
			name: "Success",
			mockSetup: func(mock sqlmock.Sqlmock) {
				// SystemStatus のモック
				mock.ExpectQuery("SELECT (.+) FROM system_statuses").
					WillReturnRows(sqlmock.NewRows([]string{"system_status_key", "login_permission", "system_state"}).AddRow("001", "1", "1"))

				// DateInfo のモック
				mock.ExpectQuery("SELECT (.+) FROM date_infos").
					WillReturnRows(sqlmock.NewRows([]string{"date_key", "prev_business_day1", "the_day", "next_business_day1", "stock_delivery_date"}).
						AddRow("001", "20231031", "20231101", "20231102", "20231106"))

				// CallPrices のモック
				mock.ExpectQuery("SELECT (.+) FROM call_prices").
					WillReturnRows(sqlmock.NewRows([]string{"unit_number", "apply_date",
						"price1", "price2", "price3", "price4", "price5", "price6", "price7", "price8", "price9", "price10",
						"price11", "price12", "price13", "price14", "price15", "price16", "price17", "price18", "price19", "price20",
						"unit_price1", "unit_price2", "unit_price3", "unit_price4", "unit_price5", "unit_price6", "unit_price7", "unit_price8", "unit_price9", "unit_price10",
						"unit_price11", "unit_price12", "unit_price13", "unit_price14", "unit_price15", "unit_price16", "unit_price17", "unit_price18", "unit_price19", "unit_price20",
						"decimal1", "decimal2", "decimal3", "decimal4", "decimal5", "decimal6", "decimal7", "decimal8", "decimal9", "decimal10",
						"decimal11", "decimal12", "decimal13", "decimal14", "decimal15", "decimal16", "decimal17", "decimal18", "decimal19", "decimal20"}).
						AddRow(101, "20231101", 3000.0, 5000.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
						AddRow(102, "20231101", 10000.0, 30000.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0))
				// IssueMasters のモック
				mock.ExpectQuery("SELECT (.+) FROM issue_masters").
					WillReturnRows(sqlmock.NewRows([]string{"issue_code", "issue_name", "trading_unit", "tokutei_f"}).
						AddRow("7203", "Toyota", 100, "1").
						AddRow("8306", "Mitsubishi UFJ", 100, "1"))

				// IssueMarketMasters のモック
				mock.ExpectQuery("SELECT (.+) FROM issue_market_masters").
					WillReturnRows(sqlmock.NewRows([]string{"issue_code", "market_code", "price_range_min", "price_range_max", "sinyou_c", "previous_close", "issue_kubun_c", "zyouzyou_kubun", "call_price_unit_number", "call_price_unit_number_yoku"}).
						AddRow("7203", "00", 100.0, 10000.0, "1", 6000.0, "01", "01", "101", "101").
						AddRow("8306", "00", 10.0, 1000.0, "1", 800.0, "01", "01", "102", "102"))

				// IssueMarketRegulations のモック
				mock.ExpectQuery("SELECT (.+) FROM issue_market_regulations").
					WillReturnRows(sqlmock.NewRows([]string{"issue_code", "listed_market", "stop_kubun", "genbutu_urituke", "seido_sinyou_sinki_uritate", "ippan_sinyou_sinki_uritate", "sinyou_syutyu_kubun"}).
						AddRow("7203", "00", "0", "0", "1", "1", "0").
						AddRow("8306", "00", "0", "0", "1", "1", "0"))

				// OperationStatusKabu のモック
				mock.ExpectQuery("SELECT (.+) FROM operation_statuses_kabu").
					WillReturnRows(sqlmock.NewRows([]string{"listed_market", "unit", "status"}).
						AddRow("00", "01", "1").
						AddRow("01", "01", "1"))

			},
			expected: &domain.MasterData{ //期待する値
				SystemStatus: domain.SystemStatus{SystemStatusKey: "001", LoginPermission: "1", SystemState: "1"},
				DateInfo:     domain.DateInfo{DateKey: "001", PrevBusinessDay1: "20231031", TheDay: "20231101", NextBusinessDay1: "20231102", StockDeliveryDate: "20231106"},
				CallPriceMap: map[string]domain.CallPrice{
					"101": {UnitNumber: "101", ApplyDate: "20231101", Price1: 3000.0, Price2: 5000.0, UnitPrice1: 1.0, Decimal1: 0},   //省略
					"102": {UnitNumber: "102", ApplyDate: "20231101", Price1: 10000.0, Price2: 30000.0, UnitPrice1: 5.0, Decimal1: 0}, //省略
				},
				IssueMap: map[string]domain.IssueMaster{
					"7203": {IssueCode: "7203", IssueName: "Toyota", TradingUnit: 100, TokuteiF: "1"},
					"8306": {IssueCode: "8306", IssueName: "Mitsubishi UFJ", TradingUnit: 100, TokuteiF: "1"},
				},
				IssueMarketMap: map[string]map[string]domain.IssueMarketMaster{
					"7203": {
						"00": {IssueCode: "7203", MarketCode: "00", PriceRangeMin: 100.0, PriceRangeMax: 10000.0, SinyouC: "1", PreviousClose: 6000.0, IssueKubunC: "01", ZyouzyouKubun: "01", CallPriceUnitNumber: "101", CallPriceUnitNumberYoku: "101"},
					},
					"8306": {
						"00": {IssueCode: "8306", MarketCode: "00", PriceRangeMin: 10.0, PriceRangeMax: 1000.0, SinyouC: "1", PreviousClose: 800.0, IssueKubunC: "01", ZyouzyouKubun: "01", CallPriceUnitNumber: "102", CallPriceUnitNumberYoku: "102"},
					},
				},
				IssueMarketRegulationMap: map[string]map[string]domain.IssueMarketRegulation{
					"7203": {
						"00": {IssueCode: "7203", ListedMarket: "00", StopKubun: "0", GenbutuUrituke: "0", SeidoSinyouSinkiUritate: "1", IppanSinyouSinkiUritate: "1", SinyouSyutyuKubun: "0"},
					},
					"8306": {
						"00": {IssueCode: "8306", ListedMarket: "00", StopKubun: "0", GenbutuUrituke: "0", SeidoSinyouSinkiUritate: "1", IppanSinyouSinkiUritate: "1", SinyouSyutyuKubun: "0"},
					},
				},
				OperationStatusKabuMap: map[string]map[string]domain.OperationStatusKabu{
					"00": {
						"01": {ListedMarket: "00", Unit: "01", Status: "1"},
					},
					"01": {
						"01": {ListedMarket: "01", Unit: "01", Status: "1"},
					},
				},
			},

			expectedErr: nil,
		},
		{
			name: "DB Error",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM system_statuses").WillReturnError(errors.New("DB error"))
			},
			expected:    nil,
			expectedErr: errors.New("DB error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			got, err := repo.GetMasterData(context.Background())

			if tt.expectedErr != nil {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("got \n%+v, want \n%+v", got, tt.expected)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
