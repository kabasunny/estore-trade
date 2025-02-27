// internal/infrastructure/persistence/signal/tests/get_latest_signals_test.go
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

func TestSignalRepository_GetLatestSignals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db) // signal. を追加

	limit := 10
	expectedSignals := []domain.Signal{
		{ID: 1, IssueCode: "7203", Side: "buy", Priority: 1, CreatedAt: time.Now()},
		{ID: 2, IssueCode: "8306", Side: "sell", Priority: 2, CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "issue_code", "side", "priority", "created_at"})
	for _, signal := range expectedSignals {
		rows.AddRow(signal.ID, signal.IssueCode, signal.Side, signal.Priority, signal.CreatedAt)
	}

	mock.ExpectQuery("^SELECT (.+) FROM signals ORDER BY created_at DESC LIMIT").WithArgs(limit).WillReturnRows(rows)

	signals, err := repo.GetLatestSignals(context.Background(), limit)
	assert.NoError(t, err)
	assert.Len(t, signals, len(expectedSignals))
	assert.Equal(t, expectedSignals, signals)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSignalRepository_GetLatestSignals_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := signal.NewSignalRepository(db) // signal. を追加
	limit := 10
	expectedError := errors.New("database error")

	mock.ExpectQuery("^SELECT (.+) FROM signals ORDER BY created_at DESC LIMIT").WithArgs(limit).WillReturnError(expectedError)

	signals, err := repo.GetLatestSignals(context.Background(), limit)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, signals)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
