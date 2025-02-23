// internal/handler/trading.go
package handler

import (
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

type TradingHandler struct {
	tradingUsecase usecase.TradingUsecase
	logger         *zap.Logger
}
