package order

import (
	"estore-trade/internal/domain"

	"gorm.io/gorm"
)

// 新しい orderRepository インスタンスを作成するコンストラクタ関数
// db: データベース接続

func NewOrderRepository(db *gorm.DB) domain.OrderRepository { // 引数の型を *gorm.DB に変更
	return &orderRepository{db: db}
}
