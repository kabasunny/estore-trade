package tachibana

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// startEventLoop は、イベント受信ループを実行する (ゴルーチン内で呼ばれる)
func (es *EventStream) startEventLoop(ctx context.Context, req *http.Request) {
	const maxRetries = 3
	const initialBackoff = 1 * time.Second
	const timeoutDuration = 10 * time.Second

	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			es.logger.Info("Stopping EventStream, context cancelled")
			return
		case <-es.stopCh:
			es.logger.Info("Stopping EventStream")
			return
		default:
			// HTTP リクエストを送信 (ここではブロック)
			resp, err := es.conn.Do(req)
			if err != nil {
				es.logger.Error("Event stream request failed", zap.Error(err))
				retryCount++
				if retryCount > maxRetries {
					es.logger.Error("Max retries reached. Stopping EventStream.", zap.Error(err))
					return
				}
				backoff := time.Duration(1<<uint(retryCount-1)) * initialBackoff
				es.logger.Error("Error in event stream processing. Retrying...", zap.Int("retryCount", retryCount), zap.Duration("backoff", backoff), zap.Error(err))
				select {
				case <-time.After(backoff):
					continue
				case <-ctx.Done():
					es.logger.Info("Stopping EventStream, context cancelled during backoff")
					return
				}
			}

			// レスポンスボディの読み込みと処理 (エラーハンドリングは processResponseBody 内で行う)
			if resp != nil && resp.StatusCode == http.StatusOK {
				if err := es.processResponseBody(resp.Body); err != nil {
					// processResponseBody内でエラーが発生した場合、ログ出力はされているが、
					// ここでは特に何もしない (リトライはしない)
				}
				resp.Body.Close()
				retryCount = 0
			} else {
				if resp != nil {
					resp.Body.Close()
				}
			}

			// 最終受信時刻から一定時間経過していたらタイムアウトと判断
			if time.Since(es.getLastReceived()) > timeoutDuration {
				es.logger.Warn("Event stream timeout. Reconnecting...")
				retryCount++ // タイムアウトもリトライカウントを増やす
				if retryCount > maxRetries {
					es.logger.Error("Max retries reached. Stopping EventStream.")
					return
				}
				backoff := time.Duration(1<<uint(retryCount-1)) * initialBackoff
				select {
				case <-time.After(backoff):
					continue // リトライ  // リトライ処理はログインを伴った方が良いかもしれない
				case <-ctx.Done():
					es.logger.Info("Stopping EventStream, context cancelled during timeout backoff")
					return
				}
			}
		}
	}
}
