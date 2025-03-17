// internal/domain/model/ranking.go
package model

// ドメインに統合する
// type Ranking struct {
// 	gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt を追加
// 	Rank         int     `gorm:"not null"`
// 	IssueCode    string  `gorm:"type:varchar(10);not null"`
// 	TradingValue float64 `gorm:"not null"`
// }

// // TableName overrides the table name used by User to `profiles`
// func (Ranking) TableName() string {
// 	return "rankings"
// }

// type RankingRepository interface {
// 	SaveRanking(ctx context.Context, ranking []Ranking) error
// 	GetLatestRanking(ctx context.Context, limit int) ([]Ranking, error)
// }
