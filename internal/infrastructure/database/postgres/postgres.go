package postgres

import (
	"database/sql"
	"estore-trade/internal/config"
	"fmt"

	// PostgreSQL driver
	"go.uber.org/zap"
)

type PostgresDB struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresDB(cfg *config.Config, logger *zap.Logger) (*PostgresDB, error) {
	// 接続文字列の作成
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	// データベースへの接続
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続確認 (ping)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return &PostgresDB{db: db, logger: logger}, nil
}

// Close closes the database connection.
func (pdb *PostgresDB) Close() error {
	return pdb.db.Close()
}

// DB returns the underlying *sql.DB instance.
func (pdb *PostgresDB) DB() *sql.DB {
	return pdb.db
}

// ここに、データベース操作のヘルパー関数を追加 (例: トランザクション処理など)。
// 例えば、以下のようなトランザクション開始関数
/*
func (pdb *PostgresDB) BeginTx(ctx context.Context) (*sql.Tx, error){
    tx, err := pdb.db.BeginTx(ctx, nil) // デフォルトの分離レベル
    if err != nil{
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    return tx, nil

}
*/
