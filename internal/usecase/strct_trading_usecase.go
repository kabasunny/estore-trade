// internal/usecase/strct_trading_usecase.go
package usecase

import (
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"

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
