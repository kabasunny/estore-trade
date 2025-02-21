package account

import (
	"database/sql"
)

// アカウントデータを扱うリポジトリ構造体
type accountRepository struct {
	db *sql.DB // データベース接続
}
