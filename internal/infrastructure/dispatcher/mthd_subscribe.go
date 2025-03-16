package dispatcher

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// Subscribe は、指定された subscriberID のチャネルを登録する
func (d *OrderEventDispatcher) Subscribe(subscriberID string, ch chan<- *domain.OrderEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.subscribers[subscriberID] = append(d.subscribers[subscriberID], ch)
	d.logger.Info("Subscriber added", zap.String("subscriberID", subscriberID))
}
