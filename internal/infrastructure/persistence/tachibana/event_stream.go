// internal/infrastructure/persistence/tachibana/event_stream.go
package tachibana

import (
	"context"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type EventStream struct {
	tachibanaClient TachibanaClient
	config          *config.Config
	logger          *zap.Logger
	eventCh         chan<- domain.OrderEvent // 修正: 送信専用チャネル
	stopCh          chan struct{}            // 停止シグナル用チャネル
	conn            *http.Client             // HTTPクライアント(長時間のポーリングに使用)
	req             *http.Request            // HTTPリクエスト
}

// NewEventStream は EventStream の新しいインスタンスを作成
func NewEventStream(client TachibanaClient, cfg *config.Config, logger *zap.Logger, eventCh chan<- domain.OrderEvent) *EventStream {
	return &EventStream{
		tachibanaClient: client,
		config:          cfg,
		logger:          logger,
		eventCh:         eventCh,
		stopCh:          make(chan struct{}),
		conn:            &http.Client{Timeout: 60 * time.Second}, // 長めのタイムアウトを設定
	}
}

// Start は EVENT I/F への接続を確立し、メッセージ受信ループを開始
func (es *EventStream) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// ログインして仮想URLを取得 (tachibanaClient.Login はセッション管理を行うように修正済み)
	requestURL, err := es.tachibanaClient.Login(ctx, "your_user_id", "your_password") // ユーザーIDとパスワードは適切に設定
	if err != nil {
		es.logger.Error("Failed to login for event stream", zap.Error(err))
		return fmt.Errorf("failed to login for event stream: %w", err)
	}

	// EVENT I/F へのリクエストURL作成
	eventURL := fmt.Sprintf("%s?p_rid=%s&p_board_no=%s&p_eno=0&p_evt_cmd=%s",
		requestURL, es.config.EventRid, es.config.EventBoardNo, es.config.EventEvtCmd)

	// HTTP GET リクエスト (Long Polling) 初回のみ
	es.req, err = http.NewRequestWithContext(ctx, http.MethodGet, eventURL, nil)
	if err != nil {
		es.logger.Error("Failed to create event stream request", zap.Error(err))
		return fmt.Errorf("failed to create event stream request: %w", err)
	}

	// メッセージ受信ループ (ゴルーチンで実行)
	for {
		select {
		case <-es.stopCh: // 停止シグナルを受け取ったら終了
			es.logger.Info("Stopping EventStream")
			return nil
		default:
			// ポーリングリクエスト送信
			resp, err := es.conn.Do(es.req) // HTTPリクエスト送信
			if err != nil {
				// ネットワークエラーやタイムアウトなど
				es.logger.Error("Event stream request failed", zap.Error(err))
				// リトライ処理 (例: 少し待ってから再接続)
				select {
				case <-time.After(5 * time.Second): // 5秒待機
					continue
				case <-es.stopCh:
					return nil // 停止指示があれば終了
				}
			}
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
			} else {
				// HTTPエラーの場合
				resp.Body.Close()
				es.logger.Error("Event stream returned non-200 status code", zap.Int("status_code", resp.StatusCode))
				// エラーに応じた処理 (例: リトライ、エラー通知など)
			}
		}
	}
}

// Stop はメッセージ受信ループを停止
func (es *EventStream) Stop() error {
	close(es.stopCh) // 停止シグナルを送信
	return nil
}

// parseEvent は、受信したメッセージをパースして domain.OrderEvent に変換
func (es *EventStream) parseEvent(message []byte) (*domain.OrderEvent, error) {
	fields := strings.Split(string(message), "^A")
	event := &domain.OrderEvent{}
	order := &domain.Order{} // 注文情報 (ECの場合)

	for _, field := range fields {
		keyValue := strings.SplitN(field, "^B", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]

		switch key {
		case "p_no": // 無視
		case "p_date":
			t, err := time.Parse("2006.01.02-15:04:05.000", value)
			if err != nil {
				es.logger.Warn("Failed to parse p_date", zap.Error(err))
				continue
			}
			event.Timestamp = t
		case "p_errno": // エラー番号
			if value != "" && value != "0" {
				errno, err := strconv.Atoi(value)
				if err != nil {
					es.logger.Warn("Failed to parse p_errno", zap.Error(err))
					continue
				}
				event.ErrNo = errno
			}
		case "p_err": // エラーメッセージ
			event.ErrMsg = value
		case "p_cmd": // コマンド (イベントタイプ)
			event.EventType = value
		case "p_ENO": // イベント番号
			eno, err := strconv.Atoi(value)
			if err != nil {
				es.logger.Warn("Failed to parse p_ENO", zap.Error(err))
				continue
			}
			event.EventNo = eno

		// EC (注文約定通知) の場合
		case "p_ON": // 注文番号
			order.ID = value
		case "p_ST": // 商品種別
			// ... (必要に応じて)
		case "p_IC": // 銘柄コード
			order.Symbol = value
		case "p_MC": // 市場コード
			// ...
		case "p_BBKB": // 売買区分
			switch value {
			case "1":
				order.Side = "sell"
			case "3":
				order.Side = "buy"
			}
		case "p_ODST": // 注文ステータス
			order.Status = value // 立花証券のステータスコード
		case "p_CRPR": // 注文価格
			price, err := strconv.ParseFloat(value, 64)
			if err == nil {
				order.Price = price
			}
		case "p_CRSR": // 注文数量
			quantity, err := strconv.Atoi(value)
			if err == nil {
				order.Quantity = quantity
			}
		// ... 他のECのフィールドも同様に処理 ...

		default: // その他の場合
			//es.logger.Warn("Unknown field in event message", zap.String("key", key)) // ログは多すぎるのでコメントアウト
		}
	}

	if event.EventType == "EC" {
		event.Order = order // ECの場合はOrder情報をセット
	}

	if event.EventType == "" {
		return nil, fmt.Errorf("event type is empty: %s", message)
	}

	return event, nil
}

// sendEvent は、パースされたイベントを usecase 層に送信
func (es *EventStream) sendEvent(event *domain.OrderEvent) {
	select {
	case es.eventCh <- *event: // チャネルに送信
	case <-es.stopCh: // 停止シグナルを受け取ったら終了
		return
	default:
		es.logger.Warn("Event channel is full, dropping event") // チャネルがフルの場合はイベントを破棄 (必要に応じてバッファリングを検討)
	}
}
