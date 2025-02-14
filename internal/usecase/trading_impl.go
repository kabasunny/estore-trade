// internal/usecase/trading_impl.go
package usecase

import (
	"context"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"go.uber.org/zap"
)

type tradingUsecase struct {
	tachibanaClient tachibana.TachibanaClient
	logger          *zap.Logger
	orderRepo       domain.OrderRepository
	accountRepo     domain.AccountRepository
	eventCh         chan *domain.OrderEvent // 修正: 双方向チャネル
}

func NewTradingUsecase(tachibanaClient tachibana.TachibanaClient, logger *zap.Logger, orderRepo domain.OrderRepository, accountRepo domain.AccountRepository) *tradingUsecase { // 修正: *TradingUsecase を返す
	return &tradingUsecase{
		tachibanaClient: tachibanaClient,
		logger:          logger,
		orderRepo:       orderRepo,
		accountRepo:     accountRepo,
		eventCh:         make(chan *domain.OrderEvent), // チャネルの初期化
	}
}

// ... (PlaceOrder, GetOrderStatus, CancelOrder メソッドは変更なし) ...
func (uc *tradingUsecase) PlaceOrder(ctx context.Context, userID, password string, order *domain.Order) (*domain.Order, error) {
	uc.logger.Info("Placing order", zap.String("user_id", userID), zap.Any("order", order))

	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password)
	if err != nil {
		uc.logger.Error("立花証券APIログインに失敗", zap.Error(err))
		return nil, err
	}

	placedOrder, err := uc.tachibanaClient.PlaceOrder(ctx, requestURL, order)
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

func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, userID, password string, orderID string) (*domain.Order, error) {
	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password)
	if err != nil {
		return nil, err
	}

	orderStatus, err := uc.tachibanaClient.GetOrderStatus(ctx, requestURL, orderID)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

func (uc *tradingUsecase) CancelOrder(ctx context.Context, userID, password string, orderID string) error {
	requestURL, err := uc.tachibanaClient.Login(ctx, userID, password)
	if err != nil {
		return err
	}

	err = uc.tachibanaClient.CancelOrder(ctx, requestURL, orderID)
	if err != nil {
		return err
	}

	return nil
}

// GetEventChannelReader は、EventStreamからイベントを受け取るためのチャネル (読み取り専用) を返す
func (uc *tradingUsecase) GetEventChannelReader() <-chan *domain.OrderEvent { // 修正
	return uc.eventCh
}

// GetEventChannelWriter は、EventStreamにイベントを送信するためのチャネル(書き込み専用)を返す
func (uc *tradingUsecase) GetEventChannelWriter() chan<- *domain.OrderEvent { // 修正
	return uc.eventCh
}

// HandleOrderEvent は、EventStreamから受け取ったイベントを処理する
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
		// ... (システムステータスに応じた処理: 例 システムが停止した場合、取引を停止するなど)
		uc.logger.Info("Received system status event", zap.Any("event", event))

	case "US": // 運用ステータス
		// ... (運用ステータスに応じた処理: 例 取引時間外になったら、注文処理を停止するなど)
		uc.logger.Info("Received operation status event", zap.Any("event", event))

	case "NS": // ニュース通知
		// ... (ニュース通知に応じた処理: 例 特定のキーワードを含むニュースを受信したら、アラートを出すなど)
		uc.logger.Info("Received news event", zap.Any("event", event))

	default:
		uc.logger.Warn("Unknown event type", zap.String("event_type", event.EventType))
	}

	return nil
}
