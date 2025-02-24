// internal/app/app.go
package app

import (
	"net/http"

	"estore-trade/internal/autotrading/auto_usecase"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

// App 構造体は、アプリケーションの依存関係を保持
type App struct {
	Config             *config.Config
	Logger             *zap.Logger
	DB                 *postgres.PostgresDB
	TachibanaClient    tachibana.TachibanaClient
	OrderRepo          domain.OrderRepository
	AccountRepo        domain.AccountRepository
	TradingUsecase     usecase.TradingUsecase
	AutoTradingUsecase auto_usecase.AutoTradingUsecase
	EventStream        *tachibana.EventStream
	HTTPServer         *http.Server
}
