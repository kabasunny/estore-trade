package tachibana

import (
	"estore-trade/internal/domain"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

// parseFD は FD メッセージをパースして event.Market に値を設定
func (es *EventStream) parseFD(event *domain.OrderEvent, key, value string) error {
	if event.Market == nil {
		event.Market = &domain.Market{}
	}

	parts := strings.Split(key, "_")
	if len(parts) != 3 {
		es.logger.Warn("Invalid p_col format", zap.String("key", key))
		return nil //parseに失敗しても、エラーにしない
	}

	//dataType := parts[0]
	rowNumber := parts[1]
	infoCode := parts[2]
	event.Market.RowNumber = rowNumber

	switch infoCode {
	case "DPP": // 現在値
		currentPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DPP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.CurrentPrice = currentPrice

	case "DPG": // 現在値前値比較
		event.Market.PriceStatus = value
	case "DV": // 出来高
		turnover, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse DV to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.Turnover = turnover

	case "BV", "GBV1", "GBV2", "GBV3", "GBV4", "GBV5", "GBV6", "GBV7", "GBV8", "GBV9", "GBV10": // 買気配数量
		bidQuantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse GBV to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.BidQuantity = bidQuantity
	case "QBP", "GBP1", "GBP2", "GBP3", "GBP4", "GBP5", "GBP6", "GBP7", "GBP8", "GBP9", "GBP10": // 買気配値
		bidPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse GBP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.BidPrice = bidPrice
	case "AV", "GAV1", "GAV2", "GAV3", "GAV4", "GAV5", "GAV6", "GAV7", "GAV8", "GAV9", "GAV10": // 売気配数量
		askQuantity, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse GAV to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.AskQuantity = askQuantity
	case "QAP", "GAP1", "GAP2", "GAP3", "GAP4", "GAP5", "GAP6", "GAP7", "GAP8", "GAP9", "GAP10": // 売気配値
		askPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse GAP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.AskPrice = askPrice
	case "DOP": // 始値
		openingPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DOP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.OpeningPrice = openingPrice
	case "DHP": // 高値
		highPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DHP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.HighPrice = highPrice
	case "DLP": // 安値
		lowPrice, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DLP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.LowPrice = lowPrice
	case "DJ": // 売買代金
		tradingVolume, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DJ to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.TradingVolume = tradingVolume
	case "DHF": // 日通し高値フラグ
		event.Market.DailyHighStatus = value
	case "DLF": // 日通し安値フラグ
		event.Market.DailyLowStatus = value
	case "QAS": //売気配値種類
		event.Market.AskQuoteType = value
	case "QBS": //買気配値種類
		event.Market.BidQuoteType = value
	case "VWAP": //VWAP
		vwap, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse VWAP to float64", zap.String("value", value), zap.Error(err))
			return nil
		}
		event.Market.VWAP = vwap
	case "DYRP": //騰落率
		ratio, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DYRP to float64", zap.String("value", value), zap.Error(err))
			return nil
		}
		event.Market.TurnoverRatio = ratio
	case "PRP": //前日終値
		close, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse PRP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.PreviousClose = close
	case "DYWP": //前日比
		change, err := strconv.ParseFloat(value, 64)
		if err != nil {
			es.logger.Warn("Failed to parse DYWP to float64", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.PreviousChange = change
	case "LISS": //所属
		event.Market.Listing = value
	case "DHP:T": //高値時刻
		highTime, err := time.Parse("15:04", value)
		if err != nil {
			es.logger.Warn("Failed to parse DHP:T", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.HighTime = highTime
	case "DLP:T": //安値時刻
		LowTime, err := time.Parse("15:04", value)
		if err != nil {
			es.logger.Warn("Failed to parse DLP:T", zap.String("value", value), zap.Error(err))
			return nil
		}
		event.Market.LowTime = LowTime
	case "DOP:T": //始値時刻
		openTime, err := time.Parse("15:04", value)
		if err != nil {
			es.logger.Warn("Failed to parse DOP:T", zap.String("value", value), zap.Error(err))
			return nil
		}
		event.Market.OpenTime = openTime
	case "DPP:T": //現在値時刻
		currentTime, err := time.Parse("15:04", value)
		if err != nil {
			es.logger.Warn("Failed to parse DPP:T", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.CurrentTime = currentTime
	case "QOV": //売-OVER
		overSell, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse QOV to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.OverSell = overSell
	case "QUV": //買-UNDER
		underBuy, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse QUV to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.Market.UnderBuy = underBuy

	default:
		es.logger.Debug("Unhandled FD infoCode", zap.String("infoCode", infoCode))
	}
	return nil
}
