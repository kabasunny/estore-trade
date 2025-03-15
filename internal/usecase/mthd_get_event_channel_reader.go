// internal/usecase/mthd_get_event_channel_reader.go
package usecase

import (
	"estore-trade/internal/domain"
)

// 注文イベントを受信するための読み取り専用チャネルを提供
func (uc *tradingUsecase) GetEventChannelReader() <-chan *domain.OrderEvent {
	// tradingUsecase の持つイベントチャネル (eventCh) の読み取り側を返す
	return uc.eventCh
}
