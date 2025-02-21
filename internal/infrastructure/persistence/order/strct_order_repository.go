package order

import (
	"database/sql"
)

type orderRepository struct {
	db *sql.DB
}
