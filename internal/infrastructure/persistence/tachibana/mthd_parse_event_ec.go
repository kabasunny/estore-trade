package tachibana

import (
	"estore-trade/internal/domain"
	"strconv"
	"time"

	"go.uber.org/zap"
)

// parseEC は EC メッセージをパースして event.Order に値を設定
func (es *EventStream) parseEC(event *domain.OrderEvent, key, value string) error {
	if event.Order == nil {
		event.Order = &domain.Order{}
	}
	switch key {
	case "p_ON": // 注文番号
		event.Order.TachibanaOrderID = value
	case "p_IC": // 銘柄コード
		event.Order.Symbol = value
	case "p_BBKB": // 売買区分
		switch value {
		case "1":
			event.Order.Side = "short"
		case "3":
			event.Order.Side = "long"
		case "5": //現渡
			event.Order.Side = "short"
		case "7": //現引
			event.Order.Side = "long"
		default:
			es.logger.Warn("Invalid p_BBKB value", zap.String("value", value))
		}
	case "p_ODST": // 注文ステータス
		event.Order.Status = value
	case "p_CRPR": //注文価格
		price, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse p_CRPR to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.Price = price
	case "p_CRSR": //注文数量
		quantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_CRSR to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.Quantity = quantity
	case "p_NT": // 通知種別
		event.Order.NotificationType = value
	case "p_ED": // 営業日
		event.Order.CreatedAt, _ = time.Parse("20060102", value) //parseに失敗しても、エラーにしない
	case "p_OON": // 親注文番号
		// event.Order にフィールドがないため、処理しない

	case "p_OT": // 注文種別
		// 使用しない

	case "p_ST": // 商品種別
		// event.Order にフィールドがないため、処理しない

	case "p_THKB": // 取引区分
		switch value {
		case "0":
			event.Order.TradeType = "spot" // 現物
		case "2":
			event.Order.TradeType = "credit_open" // 信用新規 (制度)
		case "4":
			event.Order.TradeType = "credit_close" // 信用返済 (制度)
		case "6":
			event.Order.TradeType = "credit_open" // 信用新規 (一般)
		case "8":
			event.Order.TradeType = "credit_close" // 信用返済 (一般)
		case "5": //現渡
			event.Order.TradeType = "spot" //現渡
		case "7": //現引
			event.Order.TradeType = "spot" //現引
		default:
			es.logger.Warn("Invalid p_THKB value", zap.String("value", value))
		}

	case "p_CRSJ": // 執行条件
		switch value {
		case "0":
			event.Order.ExecutionType = "" // なし
		case "2":
			event.Order.ExecutionType = "opening"
		case "4":
			event.Order.ExecutionType = "closing"
		case "6":
			event.Order.ExecutionType = "market"
		}

	case "p_CRPRKB": // 注文値段区分

	case "p_CRTKSR": // 取消数量
		canceledQuantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_CRTKSR to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.FilledQuantity -= canceledQuantity

	case "p_CREPSR": // 失効数量
		expiredQuantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_CREPSR to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.FilledQuantity -= expiredQuantity

	case "p_CREXSR": // 約定済数量
		executedQuantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_CREXSR to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.FilledQuantity = executedQuantity

	case "p_KOFG": // 繰越フラグ
	case "p_TTST": // 訂正取消ステータス
	case "p_EXST": // 約定ステータス
	case "p_LMIT": // 有効期限
		event.Order.ExpireAt, _ = time.Parse("20060102", value) //parseに失敗しても、エラーにしない
	case "p_JKK": //譲渡益課税区分
	case "p_CHNL": //チャネル
	case "p_EPRC": //失効理由コード
	case "p_EXPR": //約定値段
		executionPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse p_EXPR to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.AveragePrice = executionPrice

	case "p_EXSR": //約定数量
		executionVolume, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_EXSR to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Order.FilledQuantity = executionVolume
	case "p_EXRC": //取引所エラーコード
	case "p_EXDT": //通知日時
		event.Order.UpdatedAt, _ = time.Parse("20060102150405", value) //parseに失敗しても、エラーにしない

	case "p_IN": //銘柄名称
		// event.Order にフィールドがないため、処理しない

		// 訂正関連のフィールドは省略
	default:
		// 未知のフィールドはログ出力
		es.logger.Warn("Unknown field in EC message", zap.String("key", key))
	}
	return nil
}
