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

func TestGetDateInfoFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    *domain.DateInfo
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"date_key", "prev_business_day1", "the_day", "next_business_day1", "stock_delivery_date"}).
				AddRow("001", "20231031", "20231101", "20231102", "20231106"),
			expected: &domain.DateInfo{DateKey: "001", PrevBusinessDay1: "20231031", TheDay: "20231101", NextBusinessDay1: "20231102", StockDeliveryDate: "20231106"},
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"date_key", "prev_business_day1", "the_day", "next_business_day1", "stock_delivery_date"}),
			expected:    nil,
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
				mock.ExpectQuery("SELECT (.+) FROM date_infos").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM date_infos").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetDateInfoFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
