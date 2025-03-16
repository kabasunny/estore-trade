// internal/infrastructure/persistence/tachibana/fact_new_event_stream.go
package tachibana

import (
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher パッケージをインポート
	// OrderEventDispatcherを使うため
	"estore-trade/internal/config"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// NewEventStream は EventStream の新しいインスタンスを作成
func NewEventStream(
	client TachibanaClient,
	cfg *config.Config,
	logger *zap.Logger,
	dispatcher *dispatcher.OrderEventDispatcher) *EventStream { // OrderEventDispatcher 型に変更
	return &EventStream{
		tachibanaClient: client.(*TachibanaClientImple), // 型アサーション
		config:          cfg,                            // configをセット
		logger:          logger,
		dispatcher:      dispatcher, //  OrderEventDispatcherをセット
		stopCh:          make(chan struct{}),
		conn:            &http.Client{}, // 長めのタイムアウトを設定
		lastReceived:    time.Now(),     // 初期値として現在時刻を設定
	}
}
