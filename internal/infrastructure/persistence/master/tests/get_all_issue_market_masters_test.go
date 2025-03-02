// internal/infrastructure/persistence/master/tests/get_all_issue_market_masters_test.go
package master_test

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/master" // master パッケージをインポート
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetAllIssueMarketMastersFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    []domain.IssueMarketMaster
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"issue_code", "market_code", "price_range_min", "price_range_max", "sinyou_c", "previous_close", "issue_kubun_c", "zyouzyou_kubun", "call_price_unit_number", "call_price_unit_number_yoku"}).
				AddRow("7203", "00", 100.0, 10000.0, "1", 6000.0, "01", "01", "101", "101"). //省略
				AddRow("8306", "00", 10.0, 1000.0, "1", 800.0, "01", "01", "102", "102"),    //省略
			expected: []domain.IssueMarketMaster{
				{IssueCode: "7203", MarketCode: "00", PriceRangeMin: 100.0, PriceRangeMax: 10000.0, SinyouC: "1", PreviousClose: 6000.0, IssueKubunC: "01", ZyouzyouKubun: "01", CallPriceUnitNumber: "101", CallPriceUnitNumberYoku: "101"},
				{IssueCode: "8306", MarketCode: "00", PriceRangeMin: 10.0, PriceRangeMax: 1000.0, SinyouC: "1", PreviousClose: 800.0, IssueKubunC: "01", ZyouzyouKubun: "01", CallPriceUnitNumber: "102", CallPriceUnitNumberYoku: "102"},
			},
			expectedErr: nil,
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"issue_code", "market_code", "price_range_min", "price_range_max", "sinyou_c", "previous_close", "issue_kubun_c", "zyouzyou_kubun", "call_price_unit_number", "call_price_unit_number_yoku"}),
			expected:    nil, // nil を期待するように変更
			expectedErr: nil,
		},
		{
			name:        "DB Error",
			mockErr:     errors.New("DB error"),
			expected:    nil,
			expectedErr: errors.New("DB error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockErr != nil {
				mock.ExpectQuery("SELECT (.+) FROM issue_market_masters").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM issue_market_masters").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetAllIssueMarketMastersFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
