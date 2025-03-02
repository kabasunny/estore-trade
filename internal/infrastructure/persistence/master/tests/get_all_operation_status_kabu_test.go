// internal/infrastructure/persistence/master/tests/get_all_operation_status_kabu_test.go
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

func TestGetAllOperationStatusKabuFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    []domain.OperationStatusKabu
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"listed_market", "unit", "status"}).
				AddRow("00", "01", "1").
				AddRow("01", "01", "1"),
			expected: []domain.OperationStatusKabu{
				{ListedMarket: "00", Unit: "01", Status: "1"},
				{ListedMarket: "01", Unit: "01", Status: "1"},
			},
			expectedErr: nil,
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"listed_market", "unit", "status"}), // 空の Rows オブジェクト
			expected:    nil,                                                          // nil を期待するように変更
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
				mock.ExpectQuery("SELECT (.+) FROM operation_statuses_kabu").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM operation_statuses_kabu").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetAllOperationStatusKabuFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
