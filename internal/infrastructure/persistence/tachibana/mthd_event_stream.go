package tachibana

import (
	"context"
	"fmt"

	"estore-trade/internal/domain"
)

// ConnectEventStream は、EVENT I/F への接続を確立し、受信したイベントをチャネルに流す
func (tc *TachibanaClientImple) ConnectEventStream(ctx context.Context) (<-chan *domain.OrderEvent, error) {
	// EventStream 構造体を使うように変更
	return nil, fmt.Errorf("ConnectEventStream method should be implemented in event_stream.go")
}
