// internal/usecase/mthd_place_order.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"fmt"

	"go.uber.org/zap"
)

// APIを使用して注文を実行し、必要な事前チェックを行う
func (uc *tradingUsecase) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// 注文のログを出力
	uc.logger.Info("Placing order", zap.Any("order", order))

	// システムの稼働状態を確認
	systemStatus := uc.tachibanaClient.GetSystemStatus()
	if systemStatus.SystemState != "1" { // システム状態  0：閉局 1：開局 2：一時停止
		return nil, fmt.Errorf("system is not in service")
	}

	// 銘柄情報のチェック
	issue, ok := uc.tachibanaClient.GetIssueMaster(order.Symbol)
	if !ok {
		return nil, fmt.Errorf("invalid issue code: %s", order.Symbol)
	}

	// 売買単位のチェック 売買単位の倍数であるかを確認
	if order.Quantity%issue.TradingUnit != 0 {
		return nil, fmt.Errorf("invalid order quantity. must be multiple of %d", issue.TradingUnit)
	}

	// 呼値のチェック (tachibana パッケージの関数を使用)
	isValid, err := uc.tachibanaClient.CheckPriceIsValid(order.Symbol, order.Price, false) // 第3引数は isNextDay (当日なので false)
	if err != nil {
		return nil, fmt.Errorf("error checking price validity: %w", err)
	}
	if !isValid {
		return nil, fmt.Errorf("invalid order price: %f", order.Price)
	}

	// 立花証券APIを使用して注文を実行
	placedOrder, err := uc.tachibanaClient.PlaceOrder(ctx, order)
	if err != nil {
		uc.logger.Error("立花証券API注文実行に失敗", zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Order placed successfully", zap.String("order_id", placedOrder.ID))

	// DBに注文情報を保存 (orderRepo を使用)
	if err := uc.orderRepo.CreateOrder(ctx, placedOrder); err != nil {
		uc.logger.Error("Failed to save order to DB", zap.Error(err))
		// DB保存に失敗しても、APIからの注文自体は成功しているので、ここではエラーを返さない (ロギングはする)
	}
	return placedOrder, nil
}
