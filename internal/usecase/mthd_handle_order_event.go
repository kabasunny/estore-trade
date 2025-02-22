// internal/usecase/mthd_handle_order_event.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"

	"go.uber.org/zap"
)

func (uc *tradingUsecase) HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	uc.logger.Info("Received order event", zap.Any("event", event))

	switch event.EventType {
	case "EC": // 注文約定通知
		// 注文情報を更新
		if event.Order != nil {
			if err := uc.orderRepo.UpdateOrderStatus(ctx, event.Order.ID, event.Order.Status); err != nil {
				uc.logger.Error("Failed to update order status in DB", zap.Error(err))
				return err // DB更新失敗はリトライ可能と判断しエラーを返す
			}
		}
	case "SS": // システムステータス
		uc.logger.Info("Received system status event", zap.Any("event", event))

	case "US": // 運用ステータス
		uc.logger.Info("Received operation status event", zap.Any("event", event))

	case "NS": // ニュース通知
		uc.logger.Info("Received news event", zap.Any("event", event))

	default:
		uc.logger.Warn("Unknown event type", zap.String("event_type", event.EventType))
	}

	return nil
}
