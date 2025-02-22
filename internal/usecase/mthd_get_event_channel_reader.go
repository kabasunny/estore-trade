// internal/usecase/mthd_get_event_channel_reader.go
package usecase

import (
	"estore-trade/internal/domain"
)

func (uc *tradingUsecase) GetEventChannelReader() <-chan domain.OrderEvent {
	return uc.eventCh
}
