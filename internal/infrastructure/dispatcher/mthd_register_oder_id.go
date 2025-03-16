package dispatcher

// RegisterOrderID は、TachibanaOrderID と subscriberID の対応関係を登録
func (d *OrderEventDispatcher) RegisterOrderID(tachibanaOrderID, subscriberID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.orderIDToSubscriberID[tachibanaOrderID] = subscriberID
}
