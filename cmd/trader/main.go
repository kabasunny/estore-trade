// cmd/trader/main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"estore-trade/internal/autotrading"
	"estore-trade/internal/config"
	"estore-trade/internal/handler"
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/logger/zapLogger"
	"estore-trade/internal/infrastructure/persistence"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

func main() {
	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗: %v", err)
	}

	// ロガーの初期化
	logger, err := zapLogger.NewZapLogger(cfg)
	if err != nil {
		log.Fatalf("ロガーの初期化に失敗: %v", err)
	}
	defer logger.Sync() // プログラム終了時にロガーを同期

	// データベースの接続
	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Fatal("DB接続エラー:", zap.Error(err))
		return
	}
	defer db.Close() // プログラム終了時にデータベース接続をクローズ

	// TachibanaClientの初期化
	tachibanaClient := tachibana.NewTachibanaClient(cfg, logger)

	// マスタデータダウンロード (Login の後)
	requestURL, err := tachibanaClient.Login(context.Background(), cfg)
	if err != nil {
		logger.Fatal("Failed to login to Tachibana API", zap.Error(err))
		return
	}
	if err := tachibanaClient.DownloadMasterData(context.Background(), requestURL); err != nil {
		logger.Fatal("Failed to download master data", zap.Error(err))
		return
	}
	logger.Info("Master data downloaded successfully")

	// リポジトリの初期化
	orderRepo := persistence.NewOrderRepository(db.DB())
	accountRepo := persistence.NewAccountRepository(db.DB())

	// ユースケースの初期化
	tradingUsecase := usecase.NewTradingUsecase(tachibanaClient, logger, orderRepo, accountRepo, cfg)

	// EventStreamの初期化 (書き込み専用チャネルを渡す)
	eventStream := tachibana.NewEventStream(tachibanaClient, cfg, logger, tradingUsecase.GetEventChannelWriter())
	go func() {
		// EVENT I/F からのイベントを非同期で受信・処理
		if err := eventStream.Start(); err != nil {
			logger.Error("EventStream error", zap.Error(err))
		}
	}() // ゴルーチンでイベントストリームを開始

	// AutoTradingUsecase の初期化 (tradingUsecase, autoTradingAlgorithm, logger, config, eventCh を渡す)
	//  EventStream からのイベントを処理するゴルーチンの起動: trading usecase 層から読み取り専用チャネルを取得して使用
	autoTradingAlgorithm := &autotrading.AutoTradingAlgorithm{} // 実際のアルゴリズムのインスタンスを生成
	autoTradingUsecase := autotrading.NewAutoTradingUsecase(tradingUsecase, autoTradingAlgorithm, logger, cfg, tradingUsecase.GetEventChannelReader())
	go autoTradingUsecase.Start() // 自動売買ロジックを非同期で実行

	// ハンドラーの初期化とHTTPサーバーの設定
	tradingHandler := handler.NewTradingHandler(tradingUsecase, logger)
	http.HandleFunc("/trade", tradingHandler.HandleTrade)
	// --- 修正ここから ---
	logger.Info("Starting server on port", zap.Int("port", cfg.HTTPPort))

	// シグナル処理のコンテキストを作成
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// HTTPサーバーの初期化
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort), // config からポート番号を取得
		Handler: http.DefaultServeMux,
	}
	// --- 修正ここまで ---

	go func() {
		// エンドポイント "/trade" へのリクエストを非同期で処理
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}() // ゴルーチンでサーバーを開始

	// シグナルを待つ
	<-ctx.Done()
	logger.Info("Shutting down server...")

	// EventStreamの停止
	if err := eventStream.Stop(); err != nil {
		logger.Error("Failed to stop EventStream", zap.Error(err))
	}

	// コンテキストを使ってサーバーをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exiting")
}
