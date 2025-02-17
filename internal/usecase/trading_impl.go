// internal/usecase/trading_impl.go
package usecase

import (
	"context"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"fmt"
	"math"

	"go.uber.org/zap"
)

type tradingUsecase struct {
	tachibanaClient tachibana.TachibanaClient
	logger          *zap.Logger
	orderRepo       domain.OrderRepository
	accountRepo     domain.AccountRepository
	eventCh         chan domain.OrderEvent
	config          *config.Config // configへの参照を保持
}

// NewTradingUsecase に config を追加
func NewTradingUsecase(tachibanaClient tachibana.TachibanaClient, logger *zap.Logger, orderRepo domain.OrderRepository, accountRepo domain.AccountRepository, cfg *config.Config) *tradingUsecase {
	return &tradingUsecase{
		tachibanaClient: tachibanaClient,
		logger:          logger,
		orderRepo:       orderRepo,
		accountRepo:     accountRepo,
		eventCh:         make(chan domain.OrderEvent),
		config:          cfg, // configをセット
	}
}

func (uc *tradingUsecase) PlaceOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	uc.logger.Info("Placing order", zap.Any("order", order))

	// config から ID/Password を取得
	// Login はセッション管理を行うように修正済み
	requestURL, err := uc.tachibanaClient.Login(ctx, uc.config)
	if err != nil {
		uc.logger.Error("立花証券APIログインに失敗", zap.Error(err))
		return nil, err
	}
	//(コメント解除)以下はマスタデータを使ったチェックの例
	systemStatus := uc.tachibanaClient.GetSystemStatus()
	if systemStatus.SystemState != "1" { //  仮にシステム状態が"1"なら稼働中
		return nil, fmt.Errorf("system is not in service")
	}
	// 銘柄情報のチェック
	issue, ok := uc.tachibanaClient.GetIssueMaster(order.Symbol)
	if !ok {
		return nil, fmt.Errorf("invalid issue code: %s", order.Symbol)
	}

	// 売買単位のチェック
	if order.Quantity%issue.TradingUnit != 0 {
		return nil, fmt.Errorf("invalid order quantity. must be multiple of %d", issue.TradingUnit)
	}

	// 呼値のチェック
	callPrice, ok := uc.tachibanaClient.GetCallPrice(issue.CallPriceUnitNumber)
	if !ok {
		return nil, fmt.Errorf("call price not found for unit number: %s", issue.CallPriceUnitNumber)
	}
	if !isValidPrice(order.Price, callPrice) {
		return nil, fmt.Errorf("invalid order price: %f", order.Price)
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

// isValidPrice は、注文価格が呼値の単位に従っているかをチェックする関数 (変更なし)
func isValidPrice(price float64, callPrice tachibana.CallPrice) bool {
	prices := [20]float64{
		callPrice.Price1, callPrice.Price2, callPrice.Price3, callPrice.Price4, callPrice.Price5,
		callPrice.Price6, callPrice.Price7, callPrice.Price8, callPrice.Price9, callPrice.Price10,
		callPrice.Price11, callPrice.Price12, callPrice.Price13, callPrice.Price14, callPrice.Price15,
		callPrice.Price16, callPrice.Price17, callPrice.Price18, callPrice.Price19, callPrice.Price20,
	}
	unitPrices := [20]float64{
		callPrice.UnitPrice1, callPrice.UnitPrice2, callPrice.UnitPrice3, callPrice.UnitPrice4, callPrice.UnitPrice5,
		callPrice.UnitPrice6, callPrice.UnitPrice7, callPrice.UnitPrice8, callPrice.UnitPrice9, callPrice.UnitPrice10,
		callPrice.UnitPrice11, callPrice.UnitPrice12, callPrice.UnitPrice13, callPrice.UnitPrice14, callPrice.UnitPrice15,
		callPrice.UnitPrice16, callPrice.UnitPrice17, callPrice.UnitPrice18, callPrice.UnitPrice19, callPrice.UnitPrice20,
	}

	for i := 0; i < len(prices); i++ {
		if price <= prices[i] {
			// price が unitPrice の倍数であるか確認
			remainder := math.Mod(price, unitPrices[i])
			return remainder == 0
		}
	}

	return false // ここには到達しないはずだが、念のため
}

func (uc *tradingUsecase) GetOrderStatus(ctx context.Context, orderID string) (*domain.Order, error) {
	// config から ID/Password を取得
	// Login はセッション管理を行うように修正済み
	requestURL, err := uc.tachibanaClient.Login(ctx, uc.config)
	if err != nil {
		return nil, err
	}

	orderStatus, err := uc.tachibanaClient.GetOrderStatus(ctx, requestURL, orderID)
	if err != nil {
		return nil, err
	}

	return orderStatus, nil
}

func (uc *tradingUsecase) CancelOrder(ctx context.Context, orderID string) error {
	requestURL, err := uc.tachibanaClient.Login(ctx, uc.config)
	if err != nil {
		return err
	}
	err = uc.tachibanaClient.CancelOrder(ctx, requestURL, orderID)
	if err != nil {
		return err
	}
	return nil
}

// GetEventChannelReader は、EventStreamからイベントを受け取るためのチャネル (読み取り専用) を返す (変更なし)
func (uc *tradingUsecase) GetEventChannelReader() <-chan domain.OrderEvent {
	return uc.eventCh
}

// GetEventChannelWriter は、EventStreamにイベントを送信するためのチャネル(書き込み専用)を返す (変更なし)
func (uc *tradingUsecase) GetEventChannelWriter() chan<- domain.OrderEvent {
	return uc.eventCh
}

// HandleOrderEvent は、EventStreamから受け取ったイベントを処理する (変更なし)
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
