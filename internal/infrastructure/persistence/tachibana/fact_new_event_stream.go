package tachibana

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewEventStream は EventStream の新しいインスタンスを作成
func NewEventStream(client TachibanaClient, cfg *config.Config, logger *zap.Logger, eventCh chan<- *domain.OrderEvent) *EventStream {
	return &EventStream{
		tachibanaClient: client.(*TachibanaClientImple), // 型アサーション
		config:          cfg,                            // configをセット
		logger:          logger,
		eventCh:         eventCh,
		stopCh:          make(chan struct{}),
		conn:            &http.Client{}, // 長めのタイムアウトを設定
		lastReceived:    time.Now(),     // 初期値として現在時刻を設定
	}
}
