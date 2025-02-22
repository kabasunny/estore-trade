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
	"estore-trade/internal/infrastructure/persistence/account"
	"estore-trade/internal/infrastructure/persistence/order"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

// 極力、通信量やメモリ使用量を抑える設計を心がける
func main() {

	// 1. 初期設定
	cfg, err := config.LoadConfig(".env") // 設定ファイルの読み込み
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗: %v", err)
		return
	}

	logger, err := zapLogger.NewZapLogger(cfg) // ロガーの初期化
	if err != nil {
		log.Fatalf("ロガーの初期化に失敗: %v", err)
	}
	defer logger.Sync() // プログラム終了時にロガーを同期（バッファされたログエントリを強制的にフラッシュ））

	db, err := postgres.NewPostgresDB(cfg, logger) // データベースの接続
	if err != nil {
		logger.Fatal("DB接続エラー:", zap.Error(err))
		return
	}
	defer db.Close() // プログラム終了時にデータベース接続をクローズ

	tachibanaClient := tachibana.NewTachibanaClient(cfg, logger) // TachibanaClientの初期化

	// 2. 立花証券APIへのログインとマスタデータ取得
	err = tachibanaClient.Login(context.Background(), cfg) // 立花証券APIにログインし、仮想URL（REQUEST)を取得
	if err != nil {
		logger.Fatal("API（REQUEST I/F) のログインに失敗:", zap.Error(err))
		return
	}

	// 現状　tachibanaClientc.targetIssueCodesのデータ挿入が実装されていないので、ターゲッティングできなよ～
	if err := tachibanaClient.DownloadMasterData(context.Background()); err != nil { // マスタデータダウンロード
		logger.Fatal("マスタデータのダウンロードに失敗:", zap.Error(err))
		return
	}
	logger.Info("マスタデータのダウンロードに成功")

	// 3. リポジトリの初期化
	orderRepo := order.NewOrderRepository(db.DB())       // 注文情報を管理する
	accountRepo := account.NewAccountRepository(db.DB()) // 口座情報を管理する

	// 4. ユースケースの初期化
	tradingUsecase := usecase.NewTradingUsecase(tachibanaClient, logger, orderRepo, accountRepo, cfg)

	// 5. EventStreamの初期化 (書き込み専用チャネルを渡す)
	eventStream := tachibana.NewEventStream(tachibanaClient, cfg, logger, tradingUsecase.GetEventChannelWriter())
	// エラーチャネルを追加
	errCh := make(chan error, 1) // エラーを受け取るためのチャネル (バッファサイズ 1)

	go func() {
		if err := eventStream.Start(); err != nil { // EVENT I/F からのイベントを非同期で受信・処理
			logger.Error("EventStream error", zap.Error(err))
			errCh <- err // エラーをチャネルに送信
		}
	}()

	// 6. AutoTradingUsecase の初期化
	autoTradingAlgorithm := &autotrading.AutoTradingAlgorithm{} // 実際のアルゴリズムのインスタンスを生成　未実装　最後に実装予定
	autoTradingUsecase := autotrading.NewAutoTradingUsecase(tradingUsecase, autoTradingAlgorithm, logger, cfg, tradingUsecase.GetEventChannelReader())
	go func() {
		if err := autoTradingUsecase.Start(); err != nil {
			logger.Error("AutoTradingUsecase error", zap.Error(err))
			errCh <- err // エラーをチャネルに送信
		}
	}()

	// 7. ハンドラーの初期化とHTTPサーバーの設定
	tradingHandler := handler.NewTradingHandler(tradingUsecase, logger)
	http.HandleFunc("/trade", tradingHandler.HandleTrade)
	logger.Info("Starting server on port", zap.Int("port", cfg.HTTPPort))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM) // シグナル処理のコンテキストを作成
	defer stop()

	// HTTPサーバーの初期化
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort), // config からポート番号を取得
		Handler: http.DefaultServeMux,
	}

	go func() {
		// エンドポイント "/trade" へのリクエストを非同期で処理
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}() // ゴルーチンでサーバーを開始

	// 8. シグナル処理とサーバーのシャットダウン
	select { // 複数のチャネルからのイベントを待つ
	case <-ctx.Done(): // シグナルを待つ
	case err := <-errCh: // EventStream または AutoTradingUsecase からのエラーを待つ
		logger.Error("Received error from goroutine", zap.Error(err))
	}

	logger.Info("Shutting down server...")

	if err := eventStream.Stop(); err != nil { // EventStreamの停止
		logger.Error("Failed to stop EventStream", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // コンテキストを使ってサーバーをシャットダウン
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	logger.Info("Server exiting")
}
