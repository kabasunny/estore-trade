package tachibana

import (
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// parseSSUS は SS/US メッセージをパースして event.System に値を設定
func (es *EventStream) parseSSUS(event *domain.OrderEvent, key, value string) error {
	if event.System == nil {
		event.System = &domain.SystemStatus{}
	}
	switch key {
	case "p_CT": // 情報更新時間
		event.System.UpdateDate = value
	case "p_LK": // ログイン許可区分 (SS)
		event.System.LoginPermission = value
		event.System.LoginStatus = value //念の為残す
	case "p_SS": // システムステータス (SS)
		event.System.SystemState = value
		//event.System.SystemStatus = value  //念の為残す  <-- この行を削除
	case "p_MC":
		event.System.MarketCode = value
	case "p_UC": // 運用カテゴリー (US)
		event.System.UpdateCategory = value
	case "p_UU": // 運用ユニット (US)
		event.System.UpdateUnit = value
	case "p_US": // 運用ステータス (US)
		event.System.UpdateStatus = value
	case "p_GSCD": //原資産コード
		event.System.UnderlyingAssetCode = value
	case "p_SHSB": //商品種別
		event.System.ProductType = value
	case "p_EDK": //営業日区分
		event.System.BusinessDayFlag = value
	default:
		// 未知のフィールドはログ出力
		es.logger.Warn("Unknown field in SS/US message", zap.String("key", key))
	}
	return nil
}
