// internal/infrastructure/database/postgres/postgres.go
package postgres

import (
	"database/sql"
	"estore-trade/internal/config"
	"fmt"

	// PostgreSQL ドライバを blank import（実際には使用しないが、ドライバを登録するために必要）
	// _ "github.com/lib/pq"
	"go.uber.org/zap"
)

// PostgresDB 構造体は、PostgreSQL データベースへの接続とロガーを保持
type PostgresDB struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewPostgresDB は、新しい PostgresDB インスタンスを作成
func NewPostgresDB(cfg *config.Config, logger *zap.Logger) (*PostgresDB, error) {
	// データソース名 (DSN) を作成します。
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// DSN を使用してデータベースへの接続を開く
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// 接続エラーの場合は、エラーをログに記録して返す
		logger.Error("データベース接続に失敗", zap.Error(err)) // logger が nil の場合でもエラーを出力
		return nil, fmt.Errorf("データベース接続に失敗: %w", err)
	}

	// データベース接続をチェック
	if err := db.Ping(); err != nil {
		// Ping エラーの場合は、エラーをログに記録して返す
		logger.Error("データベース接続のチェック（ping）に失敗", zap.Error(err))
		return nil, fmt.Errorf("データベース接続のチェック（ping）が失敗: %w", err)
	}

	// データベースへの接続が成功したことをログに記録
	logger.Info("データベース接続に成功")

	return &PostgresDB{db: db, logger: logger}, nil
}

// Closeメソッド は、データベース接続を適切にクリーンアップする
func (pdb *PostgresDB) Close() error {
	return pdb.db.Close()
}

// DBメソッド は、内部の sql.DB インスタンスを返し、標準の database/sql パッケージの機能（クエリの実行やトランザクションの管理など）を利用できる
func (pdb *PostgresDB) DB() *sql.DB {
	return pdb.db
}
