package tachibana

import "fmt"

// Stop はメッセージ受信ループを停止
func (es *EventStream) Stop() error {
	fmt.Println("(es *EventStream) Stop()")
	select {
	case <-es.stopCh: //既に停止している場合は、すぐにreturn
		return nil
	default: //まだ停止していない場合は、停止処理
		close(es.stopCh)
	}
	return nil
}
