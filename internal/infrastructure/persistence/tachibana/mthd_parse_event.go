package tachibana

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// parseEvent は、受信したメッセージをパースして domain.OrderEvent に変換
func (es *EventStream) parseEvent(message []byte) (*domain.OrderEvent, error) {
	// API仕様により、項目毎に"^A"で、データを区別する
	fields := strings.Split(string(message), "^A")
	event := &domain.OrderEvent{}
	order := &domain.Order{} // 注文情報 (ECの場合)

	for _, field := range fields {
		// 項目と値を"^B"で、データを区別する
		keyValue := strings.SplitN(field, "^B", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]

		// 値を"^C"で区切り、複数の値を考慮する
		values := strings.Split(value, "^C")

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
			if value == "SS" {
				// システムステータスイベントの処理
				// (例: event.SystemStatus のようなフィールドに値を設定)
				// 今回は特に処理しない
				es.logger.Info("Received system status event") // ログだけ出力
				return event, nil                              // 処理を終了
			}
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
			// values を使うように修正 (将来、複数の値が返ってくる場合に備えて)
			if len(values) > 0 {
				order.MarketCode = values[0] // 現状は最初の値だけ使用
			}
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
