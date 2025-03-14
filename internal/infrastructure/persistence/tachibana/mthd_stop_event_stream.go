package tachibana

import "fmt"

// Stop はメッセージ受信ループを停止
func (es *EventStream) Stop() error {
	close(es.stopCh) // 停止シグナルを送信
	fmt.Println("(es *EventStream) Stop()")
	return nil
}
