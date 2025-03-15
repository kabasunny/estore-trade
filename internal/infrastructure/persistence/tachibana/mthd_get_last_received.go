package tachibana

import "time"

// getLastReceived は最終受信時刻を取得する (スレッドセーフ)
func (es *EventStream) getLastReceived() time.Time {
	es.mu.Lock()
	defer es.mu.Unlock()
	return es.lastReceived
}
