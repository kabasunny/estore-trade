// internal/app/utlil_init_app.go
package app

import (
	"context"
	"fmt"
	"net/http"

	"estore-trade/internal/autotrading/auto_algorithm"
	"estore-trade/internal/autotrading/auto_usecase"
	"estore-trade/internal/config"
	"estore-trade/internal/handler"
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/logger/zapLogger"
	"estore-trade/internal/infrastructure/persistence/account"
	"estore-trade/internal/infrastructure/persistence/master"
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"
)

func InitApp() (*App, error) {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗: %w", err)
	}

	logger, err := zapLogger.NewZapLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("ロガーの初期化に失敗: %w", err)
	}

	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("DB接続エラー: %w", err)
	}

	// MasterDataRepository のインスタンスを作成
	masterDataRepo := master.NewMasterDataRepository(db.DB())
	// DBからマスタデータを取得
	md, err := masterDataRepo.GetMasterData(context.Background())
	if err != nil {
		return nil, fmt.Errorf("DBからマスターデータの取得に失敗: %w", err)
	}

	// TachibanaClient のインスタンスを作成 (masterData を渡す)
	tachibanaClient := tachibana.NewTachibanaClient(cfg, logger, md)
	if err := tachibanaClient.Login(context.Background(), cfg); err != nil {
		return nil, fmt.Errorf("APIログインに失敗: %w", err)
	}

	orderRepo := order.NewOrderRepository(db.DB())
	accountRepo := account.NewAccountRepository(db.DB())

	tradingUsecase := usecase.NewTradingUsecase(tachibanaClient, logger, orderRepo, accountRepo, cfg)

	autoTradingAlgorithm := &auto_algorithm.AutoTradingAlgorithm{}
	autoTradingUsecase := auto_usecase.NewAutoTradingUsecase(tradingUsecase, autoTradingAlgorithm, logger, cfg, tradingUsecase.GetEventChannelReader())

	eventStream := tachibana.NewEventStream(tachibanaClient, cfg, logger, tradingUsecase.GetEventChannelWriter())

	tradingHandler := handler.NewTradingHandler(tradingUsecase, logger)
	http.HandleFunc("/trade", tradingHandler.HandleTrade)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: http.DefaultServeMux,
	}

	return &App{
		Config:             cfg,
		Logger:             logger,
		DB:                 db,
		TachibanaClient:    tachibanaClient,
		OrderRepo:          orderRepo,
		AccountRepo:        accountRepo,
		MasterDataRepo:     masterDataRepo,
		TradingUsecase:     tradingUsecase,
		AutoTradingUsecase: autoTradingUsecase,
		EventStream:        eventStream,
		HTTPServer:         httpServer,
	}, nil
}
