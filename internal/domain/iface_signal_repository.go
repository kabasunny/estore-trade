// internal/domain/iface_signal_repository.go
package domain

import (
	"context"
)

// SignalRepository はシグナルデータの永続化を担当するインターフェース
type SignalRepository interface {
	SaveSignals(ctx context.Context, signals []Signal) error
	// 必要に応じて、GetSignals, GetLatestSignals などのメソッドを追加
}
