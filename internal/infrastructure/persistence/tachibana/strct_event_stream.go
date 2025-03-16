// internal/infrastructure/persistence/tachibana/strct_event_stream.go
package tachibana

import (
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher パッケージをインポート
	// OrderEventDispatcherを使うため
	"estore-trade/internal/config"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type EventStream struct {
	tachibanaClient *TachibanaClientImple
	config          *config.Config // configへの参照を保持
	logger          *zap.Logger
	dispatcher      *dispatcher.OrderEventDispatcher // OrderEventDispatcher 型に変更
	stopCh          chan struct{}                    // 停止シグナル用チャネル
	conn            *http.Client                     // HTTPクライアント(長時間のポーリングに使用)
	lastReceived    time.Time                        // 最終受信時刻
	mu              sync.Mutex                       // lastReceived へのアクセスを保護するための Mutex
}
