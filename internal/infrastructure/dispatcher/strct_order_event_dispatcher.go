package dispatcher

import (
	"estore-trade/internal/domain"
	"sync"

	"go.uber.org/zap"
)

type OrderEventDispatcher struct {
	logger                *zap.Logger
	subscribers           map[string][]chan<- *domain.OrderEvent // subscriberID -> channels
	orderIDToSubscriberID map[string]string                      // TachibanaOrderID -> AutoTradingUsecaseのsubscriberID のマップ
	mu                    sync.RWMutex                           // 排他制御用
}
