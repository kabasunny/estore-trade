// internal/infrastructure/persistence/signal/tests/get_signals_by_issue_code_test.go
package signal

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/signal"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSignalRepository_GetSignalsByIssueCode(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db) // signal. を追加

	issueCode := "7203"
	expectedSignals := []domain.Signal{
		{ID: 1, IssueCode: issueCode, Side: "buy", Priority: 1, CreatedAt: time.Now()},
		{ID: 2, IssueCode: issueCode, Side: "sell", Priority: 2, CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "issue_code", "side", "priority", "created_at"})
	for _, signal := range expectedSignals {
		rows.AddRow(signal.ID, signal.IssueCode, signal.Side, signal.Priority, signal.CreatedAt)
	}

	mock.ExpectQuery("^SELECT (.+) FROM signals WHERE issue_code = (.+) ORDER BY created_at DESC").WithArgs(issueCode).WillReturnRows(rows)

	signals, err := repo.GetSignalsByIssueCode(context.Background(), issueCode)
	assert.NoError(t, err)
	assert.Len(t, signals, len(expectedSignals))
	assert.Equal(t, expectedSignals, signals)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSignalRepository_GetSignalsByIssueCode_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db) // signal. を追加
	issueCode := "7203"
	expectedError := errors.New("database error")

	mock.ExpectQuery("^SELECT (.+) FROM signals WHERE issue_code = (.+) ORDER BY created_at DESC").WithArgs(issueCode).WillReturnError(expectedError)

	signals, err := repo.GetSignalsByIssueCode(context.Background(), issueCode)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, signals)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
