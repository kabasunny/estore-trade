package dispatcher

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

func NewOrderEventDispatcher(logger *zap.Logger) *OrderEventDispatcher {
	return &OrderEventDispatcher{
		logger:                logger,
		subscribers:           make(map[string][]chan<- *domain.OrderEvent),
		orderIDToSubscriberID: make(map[string]string), // TachibanaOrderID -> AutoTradingUsecaseのsubscriberID のマップ
	}
}
