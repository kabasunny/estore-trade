// internal/usecase/mthd_get_event_channel_writer.go
package usecase

import (
	"estore-trade/internal/domain"
)

// 注文イベントを送信するための書き込み専用チャネルを提供

func (uc *tradingUsecase) GetEventChannelWriter() chan<- domain.OrderEvent {
	// tradingUsecase の持つイベントチャネル (eventCh) の書き込み側を返す
	return uc.eventCh
}
