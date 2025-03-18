// internal/usecase/mthd_update_order_by_event.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

func (uc *tradingUsecase) UpdateOrderByEvent(ctx context.Context, event *domain.OrderEvent) error {
	// ECイベントで、かつ Order 情報が存在することを確認
	if event.EventType != "EC" || event.Order == nil {
		return fmt.Errorf("invalid event type for order update: %s", event.EventType)
	}

	// event.Order から、必要な情報を取り出す（利用は後ほど）
	orderID := event.Order.TachibanaOrderID
	//status := event.Order.Status         // 不要
	//filledQuantity := event.Order.FilledQuantity // 不要
	//averagePrice := event.Order.AveragePrice     // 不要
	orderDate := event.Order.BusinessDate // 注文日 (営業日) を取得

	// TachibanaClient を使って、常に最新の注文情報を取得
	apiOrder, err := uc.tachibanaClient.GetOrderStatus(ctx, orderID, orderDate)
	if err != nil {
		return fmt.Errorf("failed to get order status from Tachibana API: %w", err)
	}

	// API から注文情報が取得できない場合はエラー
	if apiOrder == nil {
		return fmt.Errorf("order not found in Tachibana API: %s", orderID)
	}

	// OrderRepository を使って、DBから既存の注文情報を取得
	existingOrder, err := uc.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order from DB: %w", err)
	}

	// 既存の注文情報がない場合は、API から取得した情報で新しい Order オブジェクトを作成しDBに保存
	if existingOrder == nil {
		// API から取得した情報で新しい Order オブジェクトを作成 (DB に保存するため)
		newOrder := &domain.Order{
			//DB管理のID
			UUID: apiOrder.TachibanaOrderID, // UUID を生成
			//ここから下は、APIから取得した値で更新
			Symbol:           apiOrder.Symbol,
			Side:             apiOrder.Side,
			OrderType:        apiOrder.OrderType,
			Quantity:         apiOrder.Quantity,
			Price:            apiOrder.Price,
			Status:           apiOrder.Status,
			FilledQuantity:   apiOrder.FilledQuantity,
			AveragePrice:     apiOrder.AveragePrice,
			TachibanaOrderID: apiOrder.TachibanaOrderID, //TachibanaOrderIDを設定
			Commission:       apiOrder.Commission,
			ExpireAt:         apiOrder.ExpireAt,
		}

		// 新しい注文情報を DB に保存
		if err := uc.orderRepo.CreateOrder(ctx, newOrder); err != nil {
			uc.logger.Error("Failed to create new order in DB", zap.String("orderID", orderID), zap.Error(err))
			return fmt.Errorf("failed to create new order in DB: %w", err)
		}
		uc.logger.Info("New order created in DB based on Tachibana API response", zap.String("orderID", orderID))
	} else { //既存の情報がある場合
		// API から取得した情報で existingOrder を更新
		existingOrder.Status = apiOrder.Status                 // ステータス
		existingOrder.FilledQuantity = apiOrder.FilledQuantity // 約定数量
		existingOrder.AveragePrice = apiOrder.AveragePrice     // 平均約定価格
		// 他のフィールドも必要に応じて更新

		// UpdateOrder メソッドで更新
		err = uc.orderRepo.UpdateOrder(ctx, existingOrder) //既存のOrder情報で更新
		if err != nil {
			uc.logger.Error("Failed to update order in DB", zap.String("orderID", orderID), zap.Error(err))
			return fmt.Errorf("failed to update order in DB: %w", err)
		}
		uc.logger.Info("Order status updated", zap.String("orderID", orderID), zap.String("status", existingOrder.Status))
	}
	return nil
}
