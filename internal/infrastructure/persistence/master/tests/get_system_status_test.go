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

func TestGetSystemStatusFromDB(t *testing.T) { // テスト関数名は変更しない
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	tests := []struct {
		name        string
		mockRows    *sqlmock.Rows
		mockErr     error
		expected    *domain.SystemStatus
		expectedErr error
	}{
		{
			name: "Success",
			mockRows: sqlmock.NewRows([]string{"system_status_key", "login_permission", "system_state"}).
				AddRow("001", "1", "1"),
			expected: &domain.SystemStatus{SystemStatusKey: "001", LoginPermission: "1", SystemState: "1"},
		},
		{
			name:        "No Rows",
			mockRows:    sqlmock.NewRows([]string{"system_status_key", "login_permission", "system_state"}), // 空の結果セット
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
			if tt.mockErr != nil { //エラーの時
				mock.ExpectQuery("SELECT (.+) FROM system_statuses").WillReturnError(tt.mockErr)
			} else { //正常の時
				mock.ExpectQuery("SELECT (.+) FROM system_statuses").WillReturnRows(tt.mockRows)
			}
			got, err := master.GetSystemStatusFromDB(context.Background(), db) // master パッケージのヘルパー関数を呼び出す

			if tt.expectedErr != nil { //期待するエラーがある時
				if err == nil { //エラーがない時
					t.Errorf("expected error, got nil")
				}
			} else { //期待するエラーがない時
				if err != nil { //エラーがある時
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
