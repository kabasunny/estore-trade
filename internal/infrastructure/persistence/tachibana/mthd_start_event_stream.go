// internal/infrastructure/persistence/tachibana/mthd_start_event_stream.go
package tachibana

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Start は EVENT I/F への接続を確立し、メッセージ受信ループを開始
func (es *EventStream) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ログインして仮想URLを取得 (tachibanaClient.Login はセッション管理を行うように修正済み)
	//err := es.tachibanaClient.Login(ctx, es.config) // ログイン
	err := es.tachibanaClient.Login(ctx, nil) // ログイン 引数を削除
	if err != nil {
		es.logger.Error("Failed to login for event stream", zap.Error(err))
		return fmt.Errorf("failed to login for event stream: %w", err)
	}

	baseEventURL, _ := es.tachibanaClient.GetEventURL()

	// EVENT I/F へのリクエストURL作成
	eventURL := fmt.Sprintf("%s?p_rid=%s&p_board_no=%s&p_eno=0&p_evt_cmd=%s",
		baseEventURL, es.config.EventRid, es.config.EventBoardNo, es.config.EventEvtCmd)

	// HTTP GET リクエスト (Long Polling)　を関数化
	sendAndProcessRequest := func() error {
		// HTTP GET リクエスト (Long Polling)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, eventURL, nil)
		if err != nil {
			es.logger.Error("Failed to create event stream request", zap.Error(err))
			return fmt.Errorf("failed to create event stream request: %w", err)
		}

		// ポーリングリクエスト送信
		resp, err := es.conn.Do(req)
		if err != nil {
			es.logger.Error("Event stream request failed", zap.Error(err))
			return fmt.Errorf("event stream request failed: %w", err)
		}
		defer resp.Body.Close() //必ずクローズ

		// 正常なレスポンスの場合
		if resp.StatusCode == http.StatusOK {
			// レスポンスボディの読み込みと処理
			if err := es.processResponseBody(resp); err != nil {
				return err // processResponseBody内でエラーを返せるようにする
			}

		} else {
			// HTTPエラーの場合
			es.logger.Error("Event stream returned non-200 status code", zap.Int("status_code", resp.StatusCode))
			return fmt.Errorf("event stream returned non-200 status code: %d", resp.StatusCode)
		}
		return nil
	}

	const maxRetries = 3                   // リトライ回数制限
	const initialBackoff = 1 * time.Second // 初期バックオフ時間
	retryCount := 0

	// メッセージ受信ループ (ゴルーチンで実行)
	for {
		select {
		case <-es.stopCh: // 停止シグナルを受け取ったら終了
			es.logger.Info("Stopping EventStream")
			return nil
		default:
			// リクエスト送信とレスポンス処理
			if err := sendAndProcessRequest(); err != nil {
				retryCount++
				if retryCount > maxRetries {
					es.logger.Error("Max retries reached. Stopping EventStream.", zap.Error(err))
					return fmt.Errorf("max retries reached: %w", err) // リトライ回数上限でエラー
				}

				// 指数バックオフを計算
				backoff := time.Duration(math.Pow(2, float64(retryCount))) * initialBackoff
				es.logger.Error("Error in event stream processing. Retrying...", zap.Int("retryCount", retryCount), zap.Duration("backoff", backoff), zap.Error(err))

				// リトライ処理
				select {
				case <-time.After(backoff): // 指数バックオフ時間待機
					continue
				case <-es.stopCh:
					return nil // 停止指示があれば終了
				}
			} else { // エラーなしの場合はリトライカウントをリセット
				retryCount = 0
			}
			//  ポーリング間隔を設ける（サーバーに負荷をかけすぎないように）
			time.Sleep(1 * time.Second)
		}
	}
}
