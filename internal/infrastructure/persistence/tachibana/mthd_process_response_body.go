package tachibana

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// processResponseBody はレスポンスボディを読み込み、Shift-JIS から UTF-8 への変換、parseEvent の呼び出しを行う
func (es *EventStream) processResponseBody(body io.Reader) error {
	reader := bufio.NewReaderSize(body, 8192)
	decoder := japanese.ShiftJIS.NewDecoder() // デコーダをここで作成

	for {
		message, err := reader.ReadBytes('\n') //  0x01 -> '\n' に変更
		if err != nil {
			if err == io.EOF {
				es.logger.Info("EventStream: End of response body") // Infoレベルに変更
				return nil                                          // EOF はエラーではない
			}
			es.logger.Error("Failed to read from event stream", zap.Error(err))
			return fmt.Errorf("failed to read from event stream: %w", err)
		}

		// 1. 生データの16進数ダンプ (これまで通り)
		fmt.Printf("Raw message (hex dump):\n%s\n", hex.Dump(message))

		// 2. 生データを Shift-JIS 文字列として表示(デバッグ用)
		sjisStr, _, err := transform.String(decoder, string(message))
		if err != nil {
			es.logger.Warn("Error decoding as Shift-JIS", zap.Error(err)) // Warnレベルに変更
			// Shift-JISとしてデコードできない場合は、後続処理は行わない
		} else {
			fmt.Printf("Raw message (Shift-JIS): %s\n", sjisStr)
		}

		// 3. Shift-JIS から UTF-8 への変換 (parseEvent で行っていた処理をここに移動)
		utf8Message, _, err := transform.Bytes(decoder, message)
		if err != nil {
			es.logger.Error("Failed to decode Shift-JIS to UTF-8", zap.Error(err))
			return fmt.Errorf("failed to decode Shift-JIS: %w", err) //エラーの場合は、return
		}
		//decoder.Reset() // デコーダの状態をリセット (不要かもしれないが一応)

		// 4. UTF-8 文字列の表示
		fmt.Printf("Raw message (UTF-8): %s\n", string(utf8Message))

		event, err := es.parseEvent(utf8Message) //  utf8Message -> message
		if err != nil {
			es.logger.Error("Failed to parse event stream message", zap.Error(err))
			continue
		}

		// OrderEventDispatcher を使ってイベントをディスパッチ
		es.dispatcher.Dispatch(event)

		es.setLastReceived(time.Now()) // ★ 最終受信時刻を更新 ★
	}
}
