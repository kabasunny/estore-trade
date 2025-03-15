package tachibana

import (
	"estore-trade/internal/domain"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

// parseNS は NS メッセージをパースして event.News に値を設定
func (es *EventStream) parseNS(event *domain.OrderEvent, key, value string) error {
	if event.News == nil {
		event.News = &domain.News{}
	}
	switch key {
	case "p_ID": // ニュースID
		event.News.NewsID = value
	case "p_DT": // ニュース日付
		event.News.NewsDate = value
	case "p_TM": // ニュース時刻
		event.News.NewsTime = value
	case "p_CGN": // ニュースカテゴリ数
		count, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_CGN to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.News.NewsCategoryCount = count
	case "p_CGL": // ニュースカテゴリリスト
		event.News.CategoryList = strings.Split(value, "\x03")
	case "p_GRN": // ニュースジャンル数
		// 使用しない
	case "p_GRL": // ニュースジャンルリスト
		event.News.GenreList = strings.Split(value, "\x03")
	case "p_ISN": // 関連銘柄コードリスト数
		count, err := strconv.Atoi(value)
		if err != nil {
			es.logger.Warn("Failed to parse p_ISN to int", zap.String("value", value), zap.Error(err))
			return nil //parseに失敗しても、エラーにしない
		}
		event.News.RelatedSymbolCount = count
	case "p_ISL": // 関連銘柄コードリスト
		event.News.Symbols = strings.Split(value, "\x03")
	case "p_SKF": // 使用しない
	case "p_UPD": // 使用しない
	case "p_HDL": // ニュースタイトル
		event.News.Title = value
	case "p_TX": // ニュース本文
		event.News.Body = value
	default:
		// 未知のフィールドはログ出力
		es.logger.Warn("Unknown field in NS message", zap.String("key", key))
	}
	return nil
}
