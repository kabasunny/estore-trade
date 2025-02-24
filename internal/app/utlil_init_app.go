// internal/app/app.go
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
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"
)

// InitApp はアプリケーションの初期化を行う
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

	tachibanaClient := tachibana.NewTachibanaClient(cfg, logger)
	if err := tachibanaClient.Login(context.Background(), cfg); err != nil {
		return nil, fmt.Errorf("APIログインに失敗: %w", err)
	}

	if err := tachibanaClient.DownloadMasterData(context.Background()); err != nil {
		return nil, fmt.Errorf("マスタデータダウンロードに失敗: %w", err)
	}
	logger.Info("マスタデータのダウンロードに成功")

	orderRepo := order.NewOrderRepository(db.DB())
	accountRepo := account.NewAccountRepository(db.DB())
	tradingUsecase := usecase.NewTradingUsecase(tachibanaClient, logger, orderRepo, accountRepo, cfg)

	// AutoTradingUsecase の初期化
	autoTradingAlgorithm := &auto_algorithm.AutoTradingAlgorithm{} // インスタンス生成
	autoTradingUsecase := auto_usecase.NewAutoTradingUsecase(tradingUsecase, autoTradingAlgorithm, logger, cfg, tradingUsecase.GetEventChannelReader())

	// EventStreamの初期化
	eventStream := tachibana.NewEventStream(tachibanaClient, cfg, logger, tradingUsecase.GetEventChannelWriter())

	// HTTP サーバーの初期化 (ここではまだ起動しない)
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
		TradingUsecase:     tradingUsecase,
		AutoTradingUsecase: autoTradingUsecase,
		EventStream:        eventStream,
		HTTPServer:         httpServer,
	}, nil
}
