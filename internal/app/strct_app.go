// internal/app/strct_app.go
package app

import (
	"net/http"

	"estore-trade/internal/autotrading/auto_usecase"
	"estore-trade/internal/config"
	"estore-trade/internal/domain" // 追加
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/persistence/tachibana" // 追加
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

type App struct {
	Config             *config.Config
	Logger             *zap.Logger
	DB                 *postgres.PostgresDB
	TachibanaClient    tachibana.TachibanaClient
	OrderRepo          domain.OrderRepository
	AccountRepo        domain.AccountRepository
	MasterDataRepo     domain.MasterDataRepository // 追加
	TradingUsecase     usecase.TradingUsecase
	AutoTradingUsecase auto_usecase.AutoTradingUsecase
	EventStream        *tachibana.EventStream
	HTTPServer         *http.Server
}
