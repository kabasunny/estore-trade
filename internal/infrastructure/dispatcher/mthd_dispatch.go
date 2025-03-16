// internal/infrastructure/dispatcher/dispatcher.go
package dispatcher

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// Dispatch は、受信したイベントを適切な subscriber に振り分ける
func (d *OrderEventDispatcher) Dispatch(event *domain.OrderEvent) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// イベントのコピーを作成
	eventCopy := *event

	// イベントの種類に応じて振り分け
	switch event.EventType {
	case "KP":
		// "system" 購読者に配信
		if subscribers, ok := d.subscribers["system"]; ok {
			for _, ch := range subscribers {
				go func(ch chan<- *domain.OrderEvent, event *domain.OrderEvent) {
					select {
					case ch <- event:
					default:
						d.logger.Warn("Subscriber channel full, dropping event", zap.String("subscriberID", "system"))
					}
				}(ch, &eventCopy)
			}
		}
	case "EC":
		// TachibanaOrderID がある場合、対応する subscriberID を取得し、その購読者に配信
		if event.Order != nil && event.Order.TachibanaOrderID != "" {
			subscriberID, ok := d.orderIDToSubscriberID[event.Order.TachibanaOrderID] // マップから取得
			if !ok {
				//対応関係がない場合は、ログ出力
				d.logger.Warn("TachibanaOrderID not found in dispatcher map", zap.String("TachibanaOrderID", event.Order.TachibanaOrderID))
				return
			}
			if subscribers, ok := d.subscribers[subscriberID]; ok { //購読リストから取得
				for _, ch := range subscribers {
					go func(ch chan<- *domain.OrderEvent, event *domain.OrderEvent) {
						select {
						case ch <- event:
						default:
							d.logger.Warn("Subscriber channel full, dropping event", zap.String("subscriberID", subscriberID))
						}
					}(ch, &eventCopy)
				}
			}
		}
	default:
		// 必要に応じて、他のイベントタイプの処理を追加
	}
}
