package tachibana

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
)

// Start は EVENT I/F への接続を確立し、メッセージ受信ループを開始
func (es *EventStream) Start() error {
	fmt.Println("(es *EventStream) Start()")

	baseEventURL, err := es.tachibanaClient.GetEventURL() //event_stream.goで定義したinterfaceを利用
	if err != nil {
		return fmt.Errorf("failed to get event URL: %w", err)
	}

	// URL を構築
	eventURL, err := url.Parse(baseEventURL)
	if err != nil {
		es.logger.Error("Failed to parse event base URL", zap.Error(err))
		return fmt.Errorf("failed to parse event base URL: %w", err)
	}

	// クエリパラメータを設定
	values := url.Values{}
	values.Add("p_rid", es.config.EventRid)
	values.Add("p_board_no", es.config.EventBoardNo)
	values.Add("p_eno", "0") //p_enoは固定値
	values.Add("p_evt_cmd", es.config.EventEvtCmd)
	eventURL.RawQuery = values.Encode()

	es.logger.Info("EventStream: eventURL", zap.String("url", eventURL.String()))

	// HTTP GET リクエスト (Long Polling)　を関数化
	sendAndProcessRequest := func(ctx context.Context) error { // contextを受け取る
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, eventURL.String(), nil) //ctxでリクエストを作成
		if err != nil {
			es.logger.Error("Failed to create event stream request", zap.Error(err))
			return fmt.Errorf("failed to create event stream request: %w", err)
		}

		fmt.Println("EventStream: Request created successfully")

		resp, err := es.conn.Do(req) //EventStream構造体のconnを利用
		if err != nil {
			es.logger.Error("Event stream request failed", zap.Error(err))
			return fmt.Errorf("event stream request failed: %w", err)
		}
		defer resp.Body.Close() // 必ずBodyを閉じる

		fmt.Printf("EventStream: Response status code: %d\n", resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			// レスポンスボディの読み込みと処理
			if err := es.processResponseBody(resp.Body); err != nil {
				return err
			}
		} else {
			es.logger.Error("Event stream returned non-200 status code", zap.Int("status_code", resp.StatusCode))
			return fmt.Errorf("event stream returned non-200 status code: %d", resp.StatusCode)
		}
		return nil
	}

	const maxRetries = 3
	const initialBackoff = 1 * time.Second

	go func() { // ゴルーチンで実行
		retryCount := 0
		ctx, cancel := context.WithCancel(context.Background()) // context.WithCancel をここで使用
		defer cancel()                                          // ゴルーチン終了時にキャンセル

		for {
			select {
			case <-es.stopCh:
				es.logger.Info("Stopping EventStream")
				return // 停止シグナルを受け取ったら終了
			default:
				if err := sendAndProcessRequest(ctx); err != nil { //contextを渡す
					retryCount++
					if retryCount > maxRetries {
						es.logger.Error("Max retries reached. Stopping EventStream.", zap.Error(err))
						return //returnに変更
					}

					backoff := time.Duration(1<<uint(retryCount-1)) * initialBackoff // 2の(retryCount-1)乗
					es.logger.Error("Error in event stream processing. Retrying...", zap.Int("retryCount", retryCount), zap.Duration("backoff", backoff), zap.Error(err))

					select {
					case <-time.After(backoff):
						continue
					case <-es.stopCh:
						es.logger.Info("Stopping EventStream during backoff") // 停止を検知
						return
					case <-ctx.Done(): //contextがキャンセルされた
						es.logger.Info("Stopping EventStream, context cancelled")
						return

					}
				} else {
					retryCount = 0
				}
				//time.Sleep(1 * time.Second)  //正常の場合はSleepしない。
			}
		}
	}()

	return nil
}
