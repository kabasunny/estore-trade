// internal/autotrading/auto_usecase/mthd_start.go
package auto_usecase

import (
	"go.uber.org/zap"
)

func (a *autoTradingUsecase) Start() error {
	// eventChからのイベントを処理するためのゴルーチンを起動
	go func() {
		for event := range a.eventCh {
			a.logger.Info("AutoTradingUsecase: filtering", zap.String("TachibanaOrderID", event.Order.TachibanaOrderID))
			a.HandleEvent(*event)
		}
	}()
	return nil
}
