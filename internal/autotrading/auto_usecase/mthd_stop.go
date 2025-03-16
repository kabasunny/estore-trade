// internal/autotrading/auto_usecase/mthd_stop.go
package auto_usecase

import "fmt"

func (a *autoTradingUsecase) Stop() error {
	// OrderEventDispatcher から登録解除
	a.dispatcher.Unsubscribe(a.subscriberID, a.eventCh)
	fmt.Println("(a *autoTradingUsecase) Stop()")
	return nil
}
