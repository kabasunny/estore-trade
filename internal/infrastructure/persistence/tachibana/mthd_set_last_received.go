package tachibana

import "time"

// setLastReceived は最終受信時刻を更新する (スレッドセーフ)
func (es *EventStream) setLastReceived(t time.Time) {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.lastReceived = t
}
