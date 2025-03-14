package tachibana

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// sendEvent は、パースされたイベントを usecase 層に送信
func (es *EventStream) sendEvent(event *domain.OrderEvent) {

	select {
	case es.eventCh <- event: // チャネルに送信
		es.logger.Info("Event sent to channel", zap.Any("event", event)) // チャネル送信成功ログ
	case <-es.stopCh: // 停止シグナルを受け取ったら終了
		es.logger.Info("Event stop channel")
		return
	default:
		es.logger.Warn("Event channel is full, dropping event") // チャネルがフルの場合はイベントを破棄 (必要に応じてバッファリングを検討)
	}
}
