package tachibana

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Start は EVENT I/F への接続を確立し、メッセージ受信ループを開始
func (es *EventStream) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ログインして仮想URLを取得 (tachibanaClient.Login はセッション管理を行うように修正済み)
	startTime := time.Now()
	err := es.tachibanaClient.Login(ctx, es.config) // (1) ログイン
	loginDuration := time.Since(startTime)
	es.logger.Info("Login 完了", zap.Duration("duration", loginDuration), zap.Error(err)) // Login の処理時間をログ出力
	if err != nil {
		es.logger.Error("Failed to login for event stream", zap.Error(err))
		return fmt.Errorf("failed to login for event stream: %w", err)
	}

	startTime = time.Now()
	eventURL, err := es.tachibanaClient.GetEventURL() // (2) イベント URL の取得
	getUrlDuration := time.Since(startTime)
	es.logger.Info("GetEventURL 完了", zap.Duration("duration", getUrlDuration), zap.Error(err)) // GetEventURL の処理時間をログ出力
	if err != nil {
		es.logger.Error("Failed to get event URL", zap.Error(err))
		return fmt.Errorf("failed to get event URL: %w", err)
	}

	// EVENT I/F へのリクエストURL作成
	eventURL = fmt.Sprintf("%s?p_rid=%s&p_board_no=%s&p_eno=0&p_evt_cmd=%s",
		eventURL, es.config.EventRid, es.config.EventBoardNo, es.config.EventEvtCmd)

	// HTTP GET リクエスト (Long Polling) 初回のみ
	startTime = time.Now()
	es.req, err = http.NewRequestWithContext(ctx, http.MethodGet, eventURL, nil) // (3) リクエストの作成
	reqDuration := time.Since(startTime)
	es.logger.Info("CreateRequest 完了", zap.Duration("duration", reqDuration), zap.Error(err)) // CreateRequest の処理時間をログ出力
	if err != nil {
		es.logger.Error("Failed to create event stream request", zap.Error(err))
		return fmt.Errorf("failed to create event stream request: %w", err)
	}

	es.logger.Info("Starting EventStream loop") // ループ開始時
	// メッセージ受信ループ (ゴルーチンで実行)
	for {
		select {
		case <-es.stopCh: // 停止シグナルを受け取ったら終了
			es.logger.Info("Stopping EventStream")
			return nil
		default:
			es.logger.Info("Sending request...") // リクエスト送信前
			// ポーリングリクエスト送信
			resp, err := es.conn.Do(es.req) // HTTPリクエスト送信
			if err != nil {
				// ネットワークエラーやタイムアウトなど
				es.logger.Error("Event stream request failed", zap.Error(err))
				// リトライ処理 (例: 少し待ってから再接続)
				select {
				case <-time.After(5 * time.Second): // 5秒待機
					es.logger.Info("Retrying after 5 seconds...") // リトライ前
					continue
				case <-es.stopCh:
					return nil // 停止指示があれば終了
				}
			}
			es.logger.Info("Response received", zap.Int("status_code", resp.StatusCode)) // レスポンス受信後

			// 正常なレスポンスの場合
			if resp.StatusCode == http.StatusOK {
				// レスポンスボディの読み込み
				body, err := io.ReadAll(resp.Body) // io.ReadAll を使用
				resp.Body.Close()                  // Closeは必ず行う

				if err != nil {
					es.logger.Error("Failed to read event stream response", zap.Error(err))
					continue // 読み込み失敗時は次のループへ
				}
				// 受信データが空でなければ処理
				receivedData := string(body) // string型に変換
				if receivedData != "" {
					es.logger.Info("Received event stream message", zap.String("message", receivedData))
					// メッセージのパース処理 (parseEvent メソッドを呼び出す)
					event, err := es.parseEvent(body) // []byteを渡す
					if err != nil {
						es.logger.Error("Failed to parse event stream message", zap.Error(err))
						continue
					}
					// usecase層への通知 (sendEvent メソッドを呼び出す)
					es.sendEvent(event)
				}
				time.Sleep(1000 * time.Millisecond) //
			} else {
				// HTTPエラーの場合
				resp.Body.Close()
				es.logger.Error("Event stream returned non-200 status code", zap.Int("status_code", resp.StatusCode))
				// エラーに応じた処理 (例: リトライ、エラー通知など)
				// 今回は、エラーをログ出力するだけで、リトライや停止は行わない
			}
		}
	}
}
