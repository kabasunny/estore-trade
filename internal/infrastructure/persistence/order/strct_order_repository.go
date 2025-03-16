// internal/infrastructure/persistence/order/strct_order_repository.go
package order

import (
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB // *sql.DB を *gorm.DB に変更
}
