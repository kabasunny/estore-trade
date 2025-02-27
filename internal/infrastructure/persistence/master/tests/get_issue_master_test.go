// internal/infrastructure/persistence/master/tests/get_issue_master_test.go
package master

import (
	"context"
	"database/sql"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/master"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMasterDataRepository_GetIssueMaster(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db) // 修正

	issueCode := "7203"
	expectedIssueMaster := &domain.IssueMaster{
		IssueCode:   issueCode,
		IssueName:   "Toyota",
		TradingUnit: 100,
		TokuteiF:    "1",
	}

	rows := sqlmock.NewRows([]string{"issue_code", "issue_name", "trading_unit", "tokutei_f"}).
		AddRow(expectedIssueMaster.IssueCode, expectedIssueMaster.IssueName, expectedIssueMaster.TradingUnit, 1) // "1" を整数 1 に

	mock.ExpectQuery("^SELECT issue_code, issue_name, trading_unit, tokutei_f FROM issue_masters WHERE issue_code =").
		WithArgs(issueCode).
		WillReturnRows(rows)

	issueMaster, err := repo.GetIssueMaster(context.Background(), issueCode)
	assert.NoError(t, err)
	assert.Equal(t, expectedIssueMaster, issueMaster)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestMasterDataRepository_GetIssueMaster_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db) // 修正

	issueCode := "non-existent-code"

	mock.ExpectQuery("^SELECT issue_code, issue_name, trading_unit, tokutei_f FROM issue_masters WHERE issue_code =").
		WithArgs(issueCode).
		WillReturnError(sql.ErrNoRows)

	issueMaster, err := repo.GetIssueMaster(context.Background(), issueCode)

	assert.NoError(t, err)     // Not Found はエラーではない
	assert.Nil(t, issueMaster) // 見つからない場合は nil

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
