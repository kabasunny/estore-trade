// internal/infrastructure/persistence/master/tests/get_all_issue_masters_test.go
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

func TestGetAllIssueMastersFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    []domain.IssueMaster
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"issue_code", "issue_name", "trading_unit", "tokutei_f"}).
				AddRow("7203", "Toyota", 100, "1").
				AddRow("8306", "Mitsubishi UFJ", 100, "1"),
			expected: []domain.IssueMaster{
				{IssueCode: "7203", IssueName: "Toyota", TradingUnit: 100, TokuteiF: "1"},
				{IssueCode: "8306", IssueName: "Mitsubishi UFJ", TradingUnit: 100, TokuteiF: "1"},
			},
			expectedErr: nil,
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"issue_code", "issue_name", "trading_unit", "tokutei_f"}), // 空の Rows オブジェクト
			expected:    nil,                                                                                // nil スライスを期待するように変更
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
				mock.ExpectQuery("SELECT (.+) FROM issue_masters").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM issue_masters").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetAllIssueMastersFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
