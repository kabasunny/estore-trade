// internal/usecase/mthd_handle_order_event.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

// 受信した注文イベントの種類に応じて、適切な処理を実行
func (uc *tradingUsecase) HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	// 注文イベントを受信したことをログ出力
	uc.logger.Info("Received order event", zap.Any("event", event))

	// 注文イベントの種類に応じた処理を分岐
	switch event.EventType {
	case "EC": // 注文約定通知
		// 注文情報を更新
		if event.Order != nil {
			// データベース上の注文状況を更新
			if err := uc.orderRepo.UpdateOrderStatus(ctx, event.Order.ID, event.Order.Status); err != nil {
				uc.logger.Error("Failed to update order status in DB", zap.Error(err))
				// DB更新失敗はリトライ可能と判断しエラーを返す
				return err
			}
		}
	case "SS": // システムステータス
		// システムステータスイベントを受信したことをログ出力
		uc.logger.Info("Received system status event", zap.Any("event", event))

	case "US": // 運用ステータス
		// 運用ステータスイベントを受信したことをログ出力
		uc.logger.Info("Received operation status event", zap.Any("event", event))

	case "NS": // ニュース通知
		// ニュースイベントを受信したことをログ出力
		uc.logger.Info("Received news event", zap.Any("event", event))

	default:
		// 未知のイベントタイプの場合、警告ログを出力
		uc.logger.Warn("Unknown event type", zap.String("event_type", event.EventType))
	}

	return nil
}
