// internal/infrastructure/database/postgres/postgres.go
package postgres

import (
	"database/sql"
	"estore-trade/internal/config"
	"fmt"

	// _ "github.com/lib/pq" // 修正: PostgreSQL ドライバを blank import
	"go.uber.org/zap"
)

type PostgresDB struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresDB(cfg *config.Config, logger *zap.Logger) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// 接続エラーの場合は、エラーをログに記録して返す
		logger.Error("Failed to connect to database", zap.Error(err)) // logger が nil の場合でもエラーを出力
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		// Ping エラーの場合は、エラーをログに記録して返す
		logger.Error("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return &PostgresDB{db: db, logger: logger}, nil
}

func (pdb *PostgresDB) Close() error {
	return pdb.db.Close()
}

func (pdb *PostgresDB) DB() *sql.DB {
	return pdb.db
}
