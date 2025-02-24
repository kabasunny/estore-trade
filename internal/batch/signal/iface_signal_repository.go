// internal/batch/signal/repository.go
package signal

import (
	"context"
	"estore-trade/internal/domain"
)

// SignalRepository はシグナルデータの永続化を担当するインターフェース
type SignalRepository interface {
	SaveSignals(ctx context.Context, signals []domain.Signal) error
	// 必要に応じて、GetSignals, GetLatestSignals などのメソッドを追加
}
