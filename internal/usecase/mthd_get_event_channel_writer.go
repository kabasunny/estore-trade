// internal/usecase/mthd_get_event_channel_writer.go
package usecase

import (
	"estore-trade/internal/domain"
)

func (uc *tradingUsecase) GetEventChannelWriter() chan<- domain.OrderEvent {
	return uc.eventCh
}
