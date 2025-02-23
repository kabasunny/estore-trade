package auto_usecase // 変更

import (
	// algorithmパッケージをインポート
	"estore-trade/internal/autotrading/auto_algorithm"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

type autoTradingUsecase struct {
	tradingUsecase       usecase.TradingUsecase
	autoTradingAlgorithm *auto_algorithm.AutoTradingAlgorithm // 型を変更
	logger               *zap.Logger
	config               config.Config
	eventCh              <-chan domain.OrderEvent
}
