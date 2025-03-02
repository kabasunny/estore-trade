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

func TestGetAllCallPricesFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    []domain.CallPrice
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"unit_number", "apply_date", "price1", "price2", "price3", "price4", "price5", "price6", "price7", "price8", "price9", "price10", "price11", "price12", "price13", "price14", "price15", "price16", "price17", "price18", "price19", "price20", "unit_price1", "unit_price2", "unit_price3", "unit_price4", "unit_price5", "unit_price6", "unit_price7", "unit_price8", "unit_price9", "unit_price10", "unit_price11", "unit_price12", "unit_price13", "unit_price14", "unit_price15", "unit_price16", "unit_price17", "unit_price18", "unit_price19", "unit_price20", "decimal1", "decimal2", "decimal3", "decimal4", "decimal5", "decimal6", "decimal7", "decimal8", "decimal9", "decimal10", "decimal11", "decimal12", "decimal13", "decimal14", "decimal15", "decimal16", "decimal17", "decimal18", "decimal19", "decimal20"}).
				AddRow(101, "20231101", 3000.0, 5000.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
				AddRow(102, "20231101", 10000.0, 30000.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 5.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0),
			expected: []domain.CallPrice{
				{UnitNumber: "101", ApplyDate: "20231101", Price1: 3000.0, Price2: 5000.0, UnitPrice1: 1.0, Decimal1: 0, Price3: 0.0, Price4: 0.0, Price5: 0.0, Price6: 0.0, Price7: 0.0, Price8: 0.0, Price9: 0.0, Price10: 0.0, Price11: 0.0, Price12: 0.0, Price13: 0.0, Price14: 0.0, Price15: 0.0, Price16: 0.0, Price17: 0.0, Price18: 0.0, Price19: 0.0, Price20: 0.0, UnitPrice2: 0.0, UnitPrice3: 0.0, UnitPrice4: 0.0, UnitPrice5: 0.0, UnitPrice6: 0.0, UnitPrice7: 0.0, UnitPrice8: 0.0, UnitPrice9: 0.0, UnitPrice10: 0.0, UnitPrice11: 0.0, UnitPrice12: 0.0, UnitPrice13: 0.0, UnitPrice14: 0.0, UnitPrice15: 0.0, UnitPrice16: 0.0, UnitPrice17: 0.0, UnitPrice18: 0.0, UnitPrice19: 0.0, UnitPrice20: 0.0, Decimal2: 0, Decimal3: 0, Decimal4: 0, Decimal5: 0, Decimal6: 0, Decimal7: 0, Decimal8: 0, Decimal9: 0, Decimal10: 0, Decimal11: 0, Decimal12: 0, Decimal13: 0, Decimal14: 0, Decimal15: 0, Decimal16: 0, Decimal17: 0, Decimal18: 0, Decimal19: 0, Decimal20: 0},
				{UnitNumber: "102", ApplyDate: "20231101", Price1: 10000.0, Price2: 30000.0, UnitPrice1: 5.0, Decimal1: 0, Price3: 0.0, Price4: 0.0, Price5: 0.0, Price6: 0.0, Price7: 0.0, Price8: 0.0, Price9: 0.0, Price10: 0.0, Price11: 0.0, Price12: 0.0, Price13: 0.0, Price14: 0.0, Price15: 0.0, Price16: 0.0, Price17: 0.0, Price18: 0.0, Price19: 0.0, Price20: 0.0, UnitPrice2: 0.0, UnitPrice3: 0.0, UnitPrice4: 0.0, UnitPrice5: 0.0, UnitPrice6: 0.0, UnitPrice7: 0.0, UnitPrice8: 0.0, UnitPrice9: 0.0, UnitPrice10: 0.0, UnitPrice11: 0.0, UnitPrice12: 0.0, UnitPrice13: 0.0, UnitPrice14: 0.0, UnitPrice15: 0.0, UnitPrice16: 0.0, UnitPrice17: 0.0, UnitPrice18: 0.0, UnitPrice19: 0.0, UnitPrice20: 0.0, Decimal2: 0, Decimal3: 0, Decimal4: 0, Decimal5: 0, Decimal6: 0, Decimal7: 0, Decimal8: 0, Decimal9: 0, Decimal10: 0, Decimal11: 0, Decimal12: 0, Decimal13: 0, Decimal14: 0, Decimal15: 0, Decimal16: 0, Decimal17: 0, Decimal18: 0, Decimal19: 0, Decimal20: 0},
			},
			expectedErr: nil,
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"unit_number", "apply_date", "price1", "price2", "unit_price1", "decimal1"}), //省略
			expected:    nil,                                                                                                   //nilを返すように変更
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
				mock.ExpectQuery("SELECT (.+) FROM call_prices").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM call_prices").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetAllCallPricesFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
