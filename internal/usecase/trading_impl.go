package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"go.uber.org/zap"
)

type tradingUsecase struct {
	tachibanaClient tachibana.TachibanaClient //型を修正
	logger          *zap.Logger
}

func NewTradingUsecase(tachibanaClient tachibana.TachibanaClient, logger *zap.Logger) TradingUsecase {
	return &tradingUsecase{
		tachibanaClient: tachibanaClient,
		logger:          logger,
	}
}

func (uc *tradingUsecase) PlaceOrder(ctx context.Context, userID, password string, order *domain.Order) (*domain.Order, error) {
	uc.logger.Info("Placing order", zap.String("user_id", userID), zap.Any("order", order))
	// ログイン (必要であれば、セッショントークンをキャッシュするなど)
	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password) //context.Contextを渡す
	if err != nil {
		uc.logger.Error("立花証券APIログインに失敗", zap.Error(err))
		return nil, err
	}

	// 注文実行
	placedOrder, err := uc.tachibanaClient.PlaceOrder(ctx, requestURL, order)
	if err != nil {
		uc.logger.Error("立花証券API注文実行に失敗", zap.Error(err))
		return nil, err
	}
	uc.logger.Info("Order placed successfully", zap.String("order_id", placedOrder.ID))
	return placedOrder, nil
}

// GetOrderStatus retrieves the status of an order.
func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, userID, password string, orderID string) (*domain.Order, error) {
	// 1. ログイン (必要に応じてセッショントークンを再利用)
	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password)
	if err != nil {
		return nil, err
	}

	// 2. 注文状況の取得 (TachibanaClientのメソッドを呼び出す)
	orderStatus, err := uc.tachibanaClient.GetOrderStatus(ctx, requestURL, orderID) // 仮のメソッド名
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

// CancelOrder cancels an existing order.
func (uc *tradingUsecase) CancelOrder(ctx context.Context, userID, password string, orderID string) error {
	// 1. ログイン (必要に応じてセッショントークンを再利用)
	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password)
	if err != nil {
		return err
	}

	// 2. 注文のキャンセル (TachibanaClientのメソッドを呼び出す)
	err = uc.tachibanaClient.CancelOrder(ctx, requestURL, orderID) // 仮のメソッド名
	if err != nil {
		return err
	}

	return nil
}
