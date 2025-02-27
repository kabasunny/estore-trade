// internal/infrastructure/persistence/signal/tests/save_signals_test.go
package signal

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/signal"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSignalRepository_SaveSignals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db)

	signals := []domain.Signal{
		{IssueCode: "7203", Side: "buy", Priority: 1, CreatedAt: time.Now()},
		{IssueCode: "8306", Side: "sell", Priority: 2, CreatedAt: time.Now()},
	}

	// 1. テーブル作成 (CREATE TABLE IF NOT EXISTS) の期待値
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS signals (
    id SERIAL PRIMARY KEY,
    issue_code VARCHAR(10) NOT NULL,
    side VARCHAR(4) NOT NULL,
    priority INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    `
	mock.ExpectExec(regexp.QuoteMeta(createTableSQL)).WillReturnResult(sqlmock.NewResult(0, 0))

	// 2. トランザクション開始の期待値
	mock.ExpectBegin()

	// 3. データの挿入 (INSERT) の期待値 (2回呼ばれる)
	insertStatement := regexp.QuoteMeta("INSERT INTO signals(issue_code, side, priority, created_at) VALUES($1, $2, $3, $4)")
	mock.ExpectPrepare(insertStatement) // Prepare
	mock.ExpectExec(insertStatement).WithArgs(signals[0].IssueCode, signals[0].Side, signals[0].Priority, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(insertStatement).WithArgs(signals[1].IssueCode, signals[1].Side, signals[1].Priority, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(2, 1))

	// 4. コミットの期待値
	mock.ExpectCommit()

	err = repo.SaveSignals(context.Background(), signals)
	assert.NoError(t, err)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSignalRepository_SaveSignals_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db)

	signals := []domain.Signal{
		{IssueCode: "7203", Side: "buy", Priority: 1},
		{IssueCode: "8306", Side: "sell", Priority: 2},
	}

	// 1. テーブル作成でエラーを返す
	createTableSQL := regexp.QuoteMeta(`
    CREATE TABLE IF NOT EXISTS signals (
        id SERIAL PRIMARY KEY,
        issue_code VARCHAR(10) NOT NULL,
        side VARCHAR(4) NOT NULL,
        priority INTEGER NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE NOT NULL
    );
    `)
	mock.ExpectExec(createTableSQL).WillReturnError(errors.New("create table error"))

	// 2. トランザクションは開始されないはず

	err = repo.SaveSignals(context.Background(), signals)
	assert.Error(t, err)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
