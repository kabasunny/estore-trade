package auto_usecase

func (a *autoTradingUsecase) Start() error {
	// eventChからのイベントを処理するためのゴルーチンを起動
	go func() {
		for event := range a.eventCh {
			a.HandleEvent(event)
		}
	}()
	return nil
}
