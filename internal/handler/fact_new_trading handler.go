// internal/handler/trading.go
package handler

import (
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

func NewTradingHandler(tradingUsecase usecase.TradingUsecase, logger *zap.Logger) *TradingHandler {
	return &TradingHandler{
		tradingUsecase: tradingUsecase,
		logger:         logger,
	}
}
