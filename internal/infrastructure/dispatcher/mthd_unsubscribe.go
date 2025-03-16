package dispatcher

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// Unsubscribe は、指定された subscriberID のチャネルを登録解除する
func (d *OrderEventDispatcher) Unsubscribe(subscriberID string, ch chan<- *domain.OrderEvent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.subscribers[subscriberID]; !ok {
		return // 既に登録解除されている
	}

	// 指定されたチャネルを削除 (スライスの要素削除)
	for i, c := range d.subscribers[subscriberID] {
		if c == ch {
			d.subscribers[subscriberID] = append(d.subscribers[subscriberID][:i], d.subscribers[subscriberID][i+1:]...)
			d.logger.Info("Subscriber removed", zap.String("subscriberID", subscriberID))
			break
		}
	}

	// もし、subscriberID に紐づくチャネルがなくなったら、キーごと削除
	if len(d.subscribers[subscriberID]) == 0 {
		delete(d.subscribers, subscriberID)
	}
}
