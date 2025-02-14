package main

import (
	"log"
	"net/http"

	"estore-trade/internal/config"
	"estore-trade/internal/handler"
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/logger/zapLogger"
	"estore-trade/internal/infrastructure/persistence/tachibana" // 立花
	"estore-trade/internal/usecase"

	"go.uber.org/zap" // go.uber.org/zap はそのままインポート
)

// 設定の読み込み: 環境変数からAPIキーやデータベース接続情報などを読み込む
// ロギング: zap ロガーを使用して、アプリケーションの動作状況を記録する
// データベース接続: PostgreSQL データベースへの接続を確立・管理する
// 外部APIクライアント: 立花証券のAPIとの通信を行うクライアントを初期化する
// ビジネスロジック (ユースケース): 注文の発注、状態取得、キャンセルの処理を行う
// HTTPハンドラ: 外部からのHTTPリクエスト（例えば、API Gatewayからのリクエスト）を受け取り、ユースケース層の処理を呼び出し、結果をレスポンスとして返す

func main() {
	// 設定の読み込み
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatal("設定ファイルの読み込みに失敗:", err)
	}

	// ロガーの初期化
	logger, err := zapLogger.NewZapLogger(cfg) // zapロガーの初期化
	if err != nil {
		log.Fatal("ロガーの初期化に失敗:", err)
	}
	defer logger.Sync() // プログラム終了時にバッファをフラッシュ

	// データベース接続
	db, err := postgres.NewPostgresDB(cfg, logger)
	if err != nil {
		logger.Fatal("DB接続エラー:", zap.Error(err)) // zap.Error(err) を使う
		return
	}
	defer db.Close()

	// APIクライアントの初期化
	tachibanaClient := tachibana.NewTachibanaClient(cfg, logger)

	// リポジトリの初期化 (DBを使用する場合)
	// tradingRepo := postgres.NewTradingRepository(db)

	// ユースケースの初期化 (立花証券APIクライアントを注入)
	tradingUsecase := usecase.NewTradingUsecase(tachibanaClient, logger) //loggerも渡す

	// HTTPハンドラの初期化 (ユースケースを注入)
	tradingHandler := handler.NewTradingHandler(tradingUsecase, logger) //loggerも渡す

	// HTTPサーバーの設定と起動 (例: net/httpを使用)
	http.HandleFunc("/trade", tradingHandler.HandleTrade) // API Gatewayからのリクエストを受け付ける
	logger.Info("Starting server on port :8080")          //loggerで開始を記録
	log.Fatal(http.ListenAndServe(":8080", nil))
}
