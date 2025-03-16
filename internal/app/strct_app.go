// internal/app/strct_app.go
package app

import (
	"net/http"

	"estore-trade/internal/autotrading/auto_usecase"
	"estore-trade/internal/config"
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/database/postgres"
	"estore-trade/internal/infrastructure/dispatcher" // dispatcher パッケージをインポート
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"estore-trade/internal/usecase"

	"go.uber.org/zap"
)

type App struct {
	Config     *config.Config       // アプリケーション全体の設定情報（環境変数などから読み込まれる）
	Logger     *zap.Logger          // 構造化ロギングのためのロガーインスタンス (zap)
	DB         *postgres.PostgresDB // PostgreSQL データベースへの接続
	HTTPServer *http.Server         // HTTPリクエストを受け付け、処理するためのHTTPサーバー

	EventStream          *tachibana.EventStream           // 立花証券の EVENT I/F からのイベントストリームを処理するインスタンス
	TachibanaClient      tachibana.TachibanaClient        // 立花証券APIとの通信を担当するクライアント
	OrderEventDispatcher *dispatcher.OrderEventDispatcher // OrderEventDispatcher を dispatcher パッケージから参照

	OrderRepo      domain.OrderRepository      // 注文 (Order) データへの永続化操作（CRUDなど）を行うリポジトリ
	AccountRepo    domain.AccountRepository    // 取引口座 (Account) データへの永続化操作を行うリポジトリ
	MasterDataRepo domain.MasterDataRepository // マスタデータ（MasterData）への永続化操作, 取得を行うリポジトリ

	TradingUsecase     usecase.TradingUsecase          // 取引に関するビジネスロジック（注文、約定処理など）をまとめたユースケース
	AutoTradingUsecase auto_usecase.AutoTradingUsecase // 自動売買に関するビジネスロジックをまとめたユースケース
}
