// internal/infrastructure/persistence/tachibana/mthd_start_event_stream.go

package tachibana

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Start は EVENT I/F への接続を確立し、メッセージ受信ループを開始
func (es *EventStream) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ログインして仮想URLを取得
	err := es.tachibanaClient.Login(ctx, es.config)
	if err != nil {
		es.logger.Error("Failed to login for event stream", zap.Error(err))
		return fmt.Errorf("failed to login for event stream: %w", err)
	}

	eventURL, err := es.tachibanaClient.GetEventURL()
	if err != nil {
		es.logger.Error("Failed to get event URL", zap.Error(err))
		return fmt.Errorf("failed to get event URL: %w", err)
	}

	// EVENT I/F へのリクエストURL作成
	eventURL = fmt.Sprintf("%s?p_rid=%s&p_board_no=%s&p_eno=0&p_evt_cmd=%s",
		eventURL, es.config.EventRid, es.config.EventBoardNo, es.config.EventEvtCmd)

	// HTTP GET リクエスト (Long Polling) 初回のみ
	es.req, err = http.NewRequestWithContext(ctx, http.MethodGet, eventURL, nil)
	if err != nil {
		es.logger.Error("Failed to create event stream request", zap.Error(err))
		return fmt.Errorf("failed to create event stream request: %w", err)
	}

	es.logger.Info("Starting EventStream loop")

	// リトライ設定
	maxRetries := 3
	retryDelay := 2 * time.Second // テストしやすいように、設定可能にすることも検討

	for attempt := 0; attempt < maxRetries; attempt++ {
		select {
		case <-es.stopCh: // 停止シグナル
			es.logger.Info("Stopping EventStream")
			return nil
		default:
			es.logger.Info("Sending request...", zap.Int("attempt", attempt+1))
			resp, err := es.conn.Do(es.req)
			if err != nil {
				es.logger.Error("Event stream request failed", zap.Error(err), zap.Int("attempt", attempt+1))
				if attempt == maxRetries-1 {
					return fmt.Errorf("failed to connect to event stream after %d attempts: %w", maxRetries, err)
				}
				select {
				case <-time.After(retryDelay):
					es.logger.Info("Retrying...", zap.Int("attempt", attempt+1))
					continue
				case <-es.stopCh:
					return nil
				}
			}

			es.logger.Info("Response received", zap.Int("status_code", resp.StatusCode))

			if resp.StatusCode == http.StatusOK {
				// レスポンスボディの処理
				err = es.processResponseBody(resp) // respをそのまま渡す
				if err != nil {
					// クローズはprocessResponseBody内
					return fmt.Errorf("failed to process response body: %w", err)
				}
				// クローズはprocessResponseBody内
				return nil // 正常終了
			} else {
				resp.Body.Close()
				es.logger.Error("Event stream returned non-200 status code", zap.Int("status_code", resp.StatusCode), zap.Int("attempt", attempt+1))
				if attempt == maxRetries-1 {
					return fmt.Errorf("failed to connect to event stream after %d attempts, status code: %d", maxRetries, resp.StatusCode)
				}
				select {
				case <-time.After(retryDelay):
					es.logger.Info("Retrying...", zap.Int("attempt", attempt+1))
					continue
				case <-es.stopCh:
					return nil
				}
			}
		}
	}
	return nil // ここには到達しないはず (リトライ回数超過でエラーを返すため)
}
