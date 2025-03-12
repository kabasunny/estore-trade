// internal/domain/model/position.go

package model

import (
	"gorm.io/gorm"
)

type Position struct {
	gorm.Model         // ID, CreatedAt, UpdatedAt, DeletedAt を追加
	AccountID  uint    `gorm:"not null"`
	IssueCode  string  `gorm:"type:varchar(10);not null"`
	Quantity   int     `gorm:"not null"`
	Price      float64 `gorm:"not null"`
}

// TableName overrides the table name used by User to `profiles`
func (Position) TableName() string {
	return "positions"
}
