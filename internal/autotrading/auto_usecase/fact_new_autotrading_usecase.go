package auto_usecase // 変更

import (
	// algorithmパッケージをインポート

	"estore-trade/internal/autotrading/auto_algorithm"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

func NewAutoTradingUsecase(tradingUsecase usecase.TradingUsecase, autoTradingAlgorithm *auto_algorithm.AutoTradingAlgorithm, logger *zap.Logger, config *config.Config, eventCh <-chan domain.OrderEvent) AutoTradingUsecase { // 戻り値の型 usecase.AutoTradingUsecase
	return &autoTradingUsecase{
		tradingUsecase:       tradingUsecase,
		autoTradingAlgorithm: autoTradingAlgorithm,
		logger:               logger,
		config:               *config, // ポインタを値渡し
		eventCh:              eventCh,
	}
}
