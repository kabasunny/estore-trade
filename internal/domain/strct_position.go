// internal/domain/model.go
package domain

// 自動売買システムの中核となるデータ構造

// 保有株
type Position struct {
	Symbol   string
	Quantity int
	Price    float64 // 平均取得単価
	Side     string
}

// 他のエンティティ (例: Trade, Stock, etc.) を必要に応じて定義
