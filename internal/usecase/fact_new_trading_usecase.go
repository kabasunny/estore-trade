// internal/usecase/util_trading_usecase.go
package usecase

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

	"go.uber.org/zap"
)

func NewTradingUsecase(tachibanaClient tachibana.TachibanaClient, logger *zap.Logger, orderRepo domain.OrderRepository, accountRepo domain.AccountRepository, cfg *config.Config) *tradingUsecase {
	return &tradingUsecase{
		tachibanaClient: tachibanaClient,
		logger:          logger,
		orderRepo:       orderRepo,
		accountRepo:     accountRepo,
		eventCh:         make(chan *domain.OrderEvent),
		config:          cfg, // configをセット
	}
}
