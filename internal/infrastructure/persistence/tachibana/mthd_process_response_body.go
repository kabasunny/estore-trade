// internal/infrastructure/persistence/tachibana/mthd_process_response_body.go
package tachibana

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

// レスポンスボディを処理するヘルパー関数
func (es *EventStream) processResponseBody(resp *http.Response) error {
	//respをそのまま利用する
	body, err := io.ReadAll(resp.Body) // io.ReadAll を使用
	defer resp.Body.Close()

	if err != nil {
		es.logger.Error("Failed to read event stream response", zap.Error(err))
		return err //読み込み失敗時はエラーを返す
	}
	receivedData := string(body)
	if receivedData != "" {
		es.logger.Info("Received event stream message", zap.String("message", receivedData))
		event, err := es.parseEvent(body)
		if err != nil {
			es.logger.Error("Failed to parse event stream message", zap.Error(err))
			return err //parseに失敗したらエラーをかえす
		}
		es.sendEvent(event)

	}
	return nil // 成功
}
