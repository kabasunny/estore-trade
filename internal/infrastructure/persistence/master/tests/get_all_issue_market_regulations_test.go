// internal/infrastructure/persistence/master/tests/get_all_issue_market_regulations_test.go
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

func TestGetAllIssueMarketRegulationsFromDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    []domain.IssueMarketRegulation
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"issue_code", "listed_market", "stop_kubun", "genbutu_urituke", "seido_sinyou_sinki_uritate", "ippan_sinyou_sinki_uritate", "sinyou_syutyu_kubun"}).
				AddRow("7203", "00", "0", "0", "1", "1", "0").
				AddRow("8306", "00", "0", "0", "1", "1", "0"),
			expected: []domain.IssueMarketRegulation{
				{IssueCode: "7203", ListedMarket: "00", StopKubun: "0", GenbutuUrituke: "0", SeidoSinyouSinkiUritate: "1", IppanSinyouSinkiUritate: "1", SinyouSyutyuKubun: "0"},
				{IssueCode: "8306", ListedMarket: "00", StopKubun: "0", GenbutuUrituke: "0", SeidoSinyouSinkiUritate: "1", IppanSinyouSinkiUritate: "1", SinyouSyutyuKubun: "0"},
			},
			expectedErr: nil,
		},
		{
			name: "No Rows",
			mockRows: sqlmock.NewRows([]string{"issue_code", "listed_market", "stop_kubun", "genbutu_urituke",
				"seido_sinyou_sinki_uritate", "ippan_sinyou_sinki_uritate", "sinyou_syutyu_kubun"}),
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
				mock.ExpectQuery("SELECT (.+) FROM issue_market_regulations").WillReturnError(tt.mockErr)
			} else {
				mock.ExpectQuery("SELECT (.+) FROM issue_market_regulations").WillReturnRows(tt.mockRows)
			}

			got, err := master.GetAllIssueMarketRegulationsFromDB(context.Background(), db) // master パッケージのヘルパー関数

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
