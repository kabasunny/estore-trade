// internal/autotrading/auto_usecase/fact_new_autotrading_usecase.go
package auto_usecase

import (
	"fmt"

	"estore-trade/internal/autotrading/auto_algorithm"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher パッケージをインポート
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

var autoTradingUsecaseIDCounter int // AutoTradingUsecase の ID

func NewAutoTradingUsecase(
	tradingUsecase usecase.TradingUsecase,
	autoTradingAlgorithm auto_algorithm.AutoTradingAlgorithm,
	logger *zap.Logger,
	config *config.Config,
	dispatcher *dispatcher.OrderEventDispatcher, // OrderEventDispatcher を受け取る
) AutoTradingUsecase {
	autoTradingUsecaseIDCounter++
	subscriberID := fmt.Sprintf("autoTradingUsecase-%d", autoTradingUsecaseIDCounter) // AutoTradingUsecase の ID を生成
	usecase := &autoTradingUsecase{
		tradingUsecase:       tradingUsecase,
		autoTradingAlgorithm: autoTradingAlgorithm,
		logger:               logger,
		config:               *config,
		eventCh:              make(chan *domain.OrderEvent),
		subscriberID:         subscriberID, // subscriberID を設定
	}

	// OrderEventDispatcher に自身を登録
	dispatcher.Subscribe(subscriberID, usecase.eventCh)
	usecase.dispatcher = dispatcher //追加

	return usecase
}
