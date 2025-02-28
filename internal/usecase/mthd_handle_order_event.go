// internal/usecase/mthd_handle_order_event.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

// 受信した注文イベントの種類に応じて、適切な処理を実行
func (uc *tradingUsecase) HandleOrderEvent(ctx context.Context, event *domain.OrderEvent) error {
	// 注文イベントを受信したことをログ出力
	uc.logger.Info("Received order event", zap.Any("event", event))

	// 注文イベントの種類に応じた処理を分岐
	switch event.EventType {
	case "EC": // 注文約定通知
		if event.Order == nil {
			return fmt.Errorf("order event of type EC must have Order data")
		}
		// データベース上の注文情報を更新
		if err := uc.orderRepo.UpdateOrder(ctx, event.Order); err != nil {
			uc.logger.Error("Failed to update order in DB", zap.Error(err))
			return fmt.Errorf("failed to update order in DB: %w", err)
		}

		//TODO:
		// // データベース上の口座情報を更新 (仮の実装)
		// if err := uc.accountRepo.UpdateAccount(ctx, &domain.Account{}); err != nil {
		//  uc.logger.Error("Failed to update account in DB", zap.Error(err))
		// 	return fmt.Errorf("failed to update account in DB: %w", err)
		// }

	case "NS": // 新規注文受付
		if event.Order == nil {
			return fmt.Errorf("order event of type NS must have Order data")
		}
		// データベース上の注文情報を更新
		if err := uc.orderRepo.UpdateOrder(ctx, event.Order); err != nil {
			uc.logger.Error("Failed to update order in DB", zap.Error(err))
			return fmt.Errorf("failed to update order in DB: %w", err)
		}

	case "US": // 取消注文受付
		if event.Order == nil {
			return fmt.Errorf("order event of type US must have Order data")
		}
		// データベース上の注文ステータスを更新
		if err := uc.orderRepo.UpdateOrderStatus(ctx, event.Order.ID, "canceled"); err != nil {
			uc.logger.Error("Failed to update order status in DB", zap.Error(err))
			return fmt.Errorf("failed to update order status in DB: %w", err)
		}

	case "SS": // システムステータス
		// システムステータスイベントを受信したことをログ出力 (必要に応じて処理を追加)
		uc.logger.Info("Received system status event", zap.Any("event", event))

	// case "US": // 運用ステータス (ここでは "US" が重複しているので、片方は削除)
	// 	// 運用ステータスイベントを受信したことをログ出力 (必要に応じて処理を追加)
	// 	uc.logger.Info("Received operation status event", zap.Any("event", event))

	// case "NS": // ニュース通知  (NS は上で処理しているので、ここでは削除)
	// 	// ニュースイベントを受信したことをログ出力
	// 	uc.logger.Info("Received news event", zap.Any("event", event))

	default:
		// 未知のイベントタイプの場合、警告ログを出力
		uc.logger.Warn("Unknown event type", zap.String("event_type", event.EventType))
	}

	return nil
}
