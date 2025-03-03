// internal/infrastructure/persistence/tachibana/mthd_process_response_body.go
package tachibana

import (
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// processResponseBody はレスポンスボディの読み込みと処理を行う
func (es *EventStream) processResponseBody(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		es.logger.Error("Failed to read event stream response", zap.Error(err))
		return fmt.Errorf("failed to read event stream response: %w", err) //エラーを返す
	}

	// 受信データが空でなければ処理
	receivedData := string(body) // string型に変換
	if receivedData != "" {
		es.logger.Info("Received event stream message", zap.String("message", receivedData))
		// メッセージのパース処理 (parseEvent メソッドを呼び出す)
		event, err := es.parseEvent(body) // []byteを渡す
		if err != nil {
			es.logger.Error("Failed to parse event stream message", zap.Error(err))
			return nil //parseに失敗した場合は、エラーにしない
		}
		// usecase層への通知 (sendEvent メソッドを呼び出す)
		es.sendEvent(event)
	}
	return nil
}
