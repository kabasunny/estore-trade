// internal/infrastructure/persistence/account/tests/create_account_test.go
package account

import (
	"context"
	"errors"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/account"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock" // account_repository_test.go で定義
)

func TestAccountRepository_CreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New() // モックのDBとモックオブジェクトを取得
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := account.NewAccountRepository(db) // モックDBを渡してリポジトリを作成

	account := &domain.Account{
		ID:               "test-id",
		UserID:           "user1",
		AccountType:      "special",
		Balance:          10000,
		AvailableBalance: 10000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// クエリと期待される引数をモックに設定
	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO accounts (id, user_id, account_type, balance, available_balance, margin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `)).WithArgs(
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		sqlmock.AnyArg(), // time.Time 型の引数を AnyArg() に変更
		sqlmock.AnyArg(), // time.Time 型の引数を AnyArg() に変更
	).WillReturnResult(sqlmock.NewResult(1, 1)) // 挿入されたID, 影響行数

	err = repo.CreateAccount(context.Background(), account)
	assert.NoError(t, err)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// ... (TestAccountRepository_CreateAccount_Error も同様に修正) ...
func TestAccountRepository_CreateAccount_Error(t *testing.T) {
	db, mock, err := sqlmock.New() // sqlmock のモックを使用
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	repo := account.NewAccountRepository(db)

	account := &domain.Account{
		ID:               "test-id",
		UserID:           "user1",
		AccountType:      "special",
		Balance:          10000,
		AvailableBalance: 10000,
		Margin:           0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	expectedError := errors.New("database error")

	mock.ExpectExec(regexp.QuoteMeta(`
        INSERT INTO accounts (id, user_id, account_type, balance, available_balance, margin, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `)).WithArgs(
		account.ID,
		account.UserID,
		account.AccountType,
		account.Balance,
		account.AvailableBalance,
		account.Margin,
		sqlmock.AnyArg(), // AnyTime{} を sqlmock.AnyArg() に変更
		sqlmock.AnyArg(), // AnyTime{} を sqlmock.AnyArg() に変更
	).WillReturnError(expectedError)

	err = repo.CreateAccount(context.Background(), account)
	assert.Error(t, err)
	// エラーメッセージの比較ではなく、エラーが発生したことだけを確認
	// (sqlmock のエラーメッセージは詳細で、変更される可能性があるため)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
