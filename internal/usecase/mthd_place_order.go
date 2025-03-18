// internal/usecase/mthd_place_order.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

// internal/usecase/mthd_place_order.go
// PlaceOrder は、APIを使用して注文を実行し、必要な事前チェックを行う
func (uc *tradingUsecase) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// 注文のログを出力
	uc.logger.Info("Placing order", zap.Any("order", order))

	// 注文のバリデーション
	if err := order.Validate(); err != nil { //ドメイン層のバリデート
		return nil, fmt.Errorf("invalid order: %w", err)
	}

	// 銘柄情報のチェック
	issue, ok := uc.tachibanaClient.GetIssueMaster(ctx, order.Symbol)
	if !ok {
		return nil, fmt.Errorf("invalid issue code: %s", order.Symbol)
	}

	// 売買単位のチェック 売買単位の倍数であるかを確認
	if order.Quantity%issue.TradingUnit != 0 {
		return nil, fmt.Errorf("invalid order quantity. must be multiple of %d", issue.TradingUnit)
	}

	// 呼値のチェック (tachibana パッケージの関数を使用)
	// 成行注文の場合は、CheckPriceIsValidを呼び出さない
	if order.OrderType != "market" && order.OrderType != "stop" {
		isValid, err := uc.tachibanaClient.CheckPriceIsValid(ctx, order.Symbol, order.Price, false) // 第3引数は isNextDay (当日なので false)
		if err != nil {
			return nil, fmt.Errorf("error checking price validity: %w", err)
		}
		if !isValid {
			return nil, fmt.Errorf("invalid order price: %f", order.Price)
		}
	}

	// 立花証券APIを使用して注文を実行
	placedOrder, err := uc.tachibanaClient.PlaceOrder(ctx, order)
	if err != nil {
		uc.logger.Error("立花証券API注文実行に失敗", zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Order placed successfully", zap.String("order_id", placedOrder.TachibanaOrderID))

	placedOrder.UUID = placedOrder.TachibanaOrderID

	//注文直後の状態を使って、DBを更新(イベント経由)
	if err := uc.UpdateOrderByEvent(ctx, &domain.OrderEvent{EventType: "EC", Order: placedOrder}); err != nil {
		// ここではエラーを返さない (ログには記録)
		uc.logger.Error("Failed to update order in DB after placing order", zap.String("orderID", placedOrder.TachibanaOrderID), zap.Error(err))
		//return nil, fmt.Errorf("failed to update order after place order: %w", err) //エラーを返した方が、より堅牢になる
	}

	return placedOrder, nil
}
