package order

import (
	"database/sql"
	"estore-trade/internal/domain"
)

// 新しい orderRepository インスタンスを作成するコンストラクタ関数
// db: データベース接続
func NewOrderRepository(db *sql.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}
