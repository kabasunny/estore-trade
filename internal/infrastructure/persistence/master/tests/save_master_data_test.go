// internal/infrastructure/persistence/master/tests/save_master_data_test.go
package master

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/master"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestMasterDataRepository_SaveMasterData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := master.NewMasterDataRepository(db) // 修正

	masterData := &domain.MasterData{
		IssueMap: map[string]domain.IssueMaster{
			"7203": {IssueCode: "7203", IssueName: "Toyota", TradingUnit: 100, TokuteiF: "1"},
			"8306": {IssueCode: "8306", IssueName: "Mitsubishi UFJ", TradingUnit: 100, TokuteiF: "1"},
		},
	}

	// トランザクション開始の期待値
	mock.ExpectBegin()

	// テーブル作成 (CREATE TABLE IF NOT EXISTS) の期待値 (今回は省略)
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS issue_masters (
        issue_code VARCHAR(10) PRIMARY KEY,
        issue_name VARCHAR(255) NOT NULL,
        trading_unit INTEGER NOT NULL,
        tokutei_f BOOLEAN NOT NULL
    );
    `
	mock.ExpectExec(regexp.QuoteMeta(createTableSQL)).WillReturnResult(sqlmock.NewResult(0, 0))

	// 既存データ削除 (DELETE) の期待値
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM issue_masters")).WillReturnResult(sqlmock.NewResult(0, 0))

	// データ挿入 (INSERT) の期待値 (2回呼ばれる)
	insertStatement := regexp.QuoteMeta("INSERT INTO issue_masters(issue_code, issue_name, trading_unit, tokutei_f) VALUES($1, $2, $3, $4)")
	mock.ExpectPrepare(insertStatement)
	mock.ExpectExec(insertStatement).WithArgs("7203", "Toyota", 100, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec(insertStatement).WithArgs("8306", "Mitsubishi UFJ", 100, 1).WillReturnResult(sqlmock.NewResult(2, 1))

	// コミットの期待値
	mock.ExpectCommit()

	err = repo.SaveMasterData(context.Background(), masterData)
	assert.NoError(t, err)

	// モックの設定がすべて満たされたか確認
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// 必要に応じて、エラーケースのテスト (BeginTx, ExecContext, PrepareContext, Commit がエラーを返す場合) も追加
