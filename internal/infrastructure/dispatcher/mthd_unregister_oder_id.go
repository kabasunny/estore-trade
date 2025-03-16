package dispatcher

// UnregisterOrderID は、TachibanaOrderID と subscriberID の対応関係を削除
func (d *OrderEventDispatcher) UnregisterOrderID(tachibanaOrderID string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.orderIDToSubscriberID, tachibanaOrderID)
}
