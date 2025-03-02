// internal/infrastructure/persistence/master/tests/get_all_issue_codes_test.go
package master_test

import (
	"context"
	"errors"
	"estore-trade/internal/infrastructure/persistence/master"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMasterDataRepository_GetAllIssueCodes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db) // 修正

	expectedIssueCodes := []string{"1301", "1332", "1333"}

	rows := sqlmock.NewRows([]string{"issue_code"}).
		AddRow("1301").
		AddRow("1332").
		AddRow("1333")

	mock.ExpectQuery("^SELECT issue_code FROM issue_masters").WillReturnRows(rows)

	issueCodes, err := repo.GetAllIssueCodes(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedIssueCodes, issueCodes)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMasterDataRepository_GetAllIssueCodes_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db) // 修正

	expectedError := errors.New("some database error") // 任意のエラー

	mock.ExpectQuery("^SELECT issue_code FROM issue_masters").WillReturnError(expectedError)

	issueCodes, err := repo.GetAllIssueCodes(context.Background())

	assert.Error(t, err)                // エラーが発生することを期待
	assert.Equal(t, expectedError, err) // 期待されるエラーと一致するか確認
	assert.Nil(t, issueCodes)           // データは返ってこない

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
