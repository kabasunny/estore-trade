package tachibana

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type EventStream struct {
	tachibanaClient TachibanaClient
	config          *config.Config // configへの参照を保持
	logger          *zap.Logger
	eventCh         chan<- *domain.OrderEvent // 修正: 送信専用チャネル
	stopCh          chan struct{}             // 停止シグナル用チャネル
	conn            *http.Client              // HTTPクライアント(長時間のポーリングに使用)
	lastReceived    time.Time                 // 最終受信時刻
	mu              sync.Mutex                // lastReceived へのアクセスを保護するための Mutex
}
