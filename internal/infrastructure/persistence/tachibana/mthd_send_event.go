package tachibana

import (
	"estore-trade/internal/domain"
)

// sendEvent は、パースされたイベントを usecase 層に送信
func (es *EventStream) sendEvent(event *domain.OrderEvent) {
	select {
	case es.eventCh <- *event: // チャネルに送信
	case <-es.stopCh: // 停止シグナルを受け取ったら終了
		return
	default:
		es.logger.Warn("Event channel is full, dropping event") // チャネルがフルの場合はイベントを破棄 (必要に応じてバッファリングを検討)
	}
}
