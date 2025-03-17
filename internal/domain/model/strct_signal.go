// internal/domain/model/signal.go
package model

// ドメインに統合する
// type Signal struct {
// 	gorm.Model        // ID, CreatedAt, UpdatedAt, DeletedAt を追加
// 	IssueCode  string `gorm:"type:varchar(10);not null"`
// 	Side       string `gorm:"type:varchar(10);not null"` // "buy" or "sell"
// 	Priority   int    `gorm:"not null"`
// }

// // TableName overrides the table name used by User to `profiles`
// func (Signal) TableName() string {
// 	return "signals"
// }

// type SignalRepository interface {
// 	SaveSignals(ctx context.Context, signals []Signal) error
// 	GetLatestSignals(ctx context.Context) ([]Signal, error)
// 	GetSignalsByIssueCode(ctx context.Context, issueCode string) ([]Signal, error)
// }
