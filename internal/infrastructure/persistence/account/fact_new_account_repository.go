package account

import (
	"database/sql"
	"estore-trade/internal/domain"
)

// NewAccountRepository は、新しい accountRepository インスタンスを作成するコンストラクタ関数
// db: データベース接続
func NewAccountRepository(db *sql.DB) domain.AccountRepository {
	return &accountRepository{db: db} // accountRepository 構造体を初期化して返す
}
