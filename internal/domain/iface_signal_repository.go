// internal/domain/iface_signal_repository.go
package domain

import (
	"context"
)

// SignalRepository はシグナルデータの永続化を担当するインターフェース
type SignalRepository interface {
	SaveSignals(ctx context.Context, signals []Signal) error
	GetSignalsByIssueCode(ctx context.Context, issueCode string) ([]Signal, error) // 追加: 銘柄コードで検索
	GetLatestSignals(ctx context.Context, limit int) ([]Signal, error)             // 追加: 最新のシグナルを取得
}
