# estore-trade/internal/app

このパッケージは、`estore-trade` アプリケーションの初期化と依存性の管理を行う

## 概要

`app` パッケージは、アプリケーションの起動に必要な各種コンポーネント（設定、ロガー、データベース接続、APIクライアント、リポジトリ、ユースケースなど）の初期化処理をカプセル化する
`InitApp` 関数は、これらのコンポーネントを初期化し、`App` 構造体にまとめて返す
`App` 構造体は、アプリケーション全体で共有されるリソースへのアクセスを提供する

## ファイル構成

-   `strct_app.go`: `App` 構造体の定義
-   `utlil_init_app.go`: `InitApp` 関数の定義（アプリケーションの初期化処理）
-   `mthd_close.go`: `Close` メソッドの定義 (アプリケーションの終了処理)

### `strct_app.go` (`App` 構造体)

-   **役割**: アプリケーション全体で共有される依存関係（設定、ロガー、DB接続など）を保持する
-   **詳細**:
    -   `Config`: アプリケーションの設定 (`*config.Config`)
    -   `Logger`: ロガー (`*zap.Logger`)
    -   `DB`: データベース接続 (`*postgres.PostgresDB`)
    -   `TachibanaClient`: 立花証券APIクライアント (`tachibana.TachibanaClient`)
    -   `OrderRepo`: 注文リポジトリ (`domain.OrderRepository`)
    -   `AccountRepo`: 口座リポジトリ (`domain.AccountRepository`)
    -   `TradingUsecase`: 取引ユースケース (`usecase.TradingUsecase`)
    -   `AutoTradingUsecase`: 自動売買ユースケース (`auto_usecase.AutoTradingUsecase`)
    -   `EventStream`: イベントストリーム (`*tachibana.EventStream`)
    -   `HTTPServer`: HTTP サーバー (`*http.Server`)

### `utlil_init_app.go` (`InitApp` 関数)

-   **役割**: アプリケーションの初期化処理を行う
-   **詳細**:
    1.  設定ファイル (`.env`) の読み込み
    2.  ロガー (zap) の初期化
    3.  PostgreSQL データベースへの接続
    4.  TachibanaClient の初期化
    5.  立花証券APIへのログイン
    6.  マスタデータのダウンロード
    7.  リポジトリ (OrderRepository, AccountRepository) の初期化
    8.  ユースケース (TradingUsecase, AutoTradingUsecase) の初期化
    9.  EventStream の初期化
    10. HTTP サーバーの初期化 (起動はしない)
    11. 初期化されたオブジェクトを `App` 構造体にまとめて返す
    12. エラーが発生した場合は、エラーを返す

### `mthd_close.go`(`Close`メソッド)
-   **役割**: アプリケーションの終了時に必要なクリーンアップ処理
-   **詳細**:
    -   データベース接続をクローズ
    -   ログをフラッシュ

## 依存関係

このパッケージは、以下のパッケージに依存する

-   `context`: コンテキスト
-   `fmt`: フォーマット付きI/O
-   `net/http`: HTTP クライアントとサーバーの実装
-   `go.uber.org/zap`: ロギングライブラリ
-   `estore-trade/internal/config`: 設定管理
-   `estore-trade/internal/domain`: ドメインモデルとリポジトリインターフェース
-   `estore-trade/internal/handler`: HTTPハンドラ
-   `estore-trade/internal/infrastructure/database/postgres`: PostgreSQL データベースドライバ
-   `estore-trade/internal/infrastructure/logger/zapLogger`: zap ロガーの実装
-   `estore-trade/internal/infrastructure/persistence/account`: アカウントリポジトリの実装
-   `estore-trade/internal/infrastructure/persistence/order`: 注文リポジトリの実装
-   `estore-trade/internal/infrastructure/persistence/tachibana`: 立花証券APIクライアントの実装
-   `estore-trade/internal/usecase`: ユースケースの実装
-   `estore-trade/internal/autotrading/auto_algorithm`:売買アルゴリズムの定義
-   `estore-trade/internal/autotrading/auto_usecase`: 自動売買ユースケースの実装

## 特記事項

-   **依存性の注入:** `App` 構造体を通じて、アプリケーションの各コンポーネントに必要な依存関係を渡すことで、各コンポーネントのテストが容易になる
-   **エラーハンドリング:** `InitApp` 関数は、初期化処理中にエラーが発生した場合、エラーを返す。呼び出し元 (通常は `main` 関数) でエラーを適切に処理する必要がある
-   **クリーンアップ:** `App` 構造体の `Close` メソッドは、アプリケーションの終了時に呼び出され、リソース (データベース接続など) を解放する。 `main` 関数では、`defer app.Close()` を呼び出して、確実にクリーンアップ処理が実行されるようにする
- **設定の外部化**: 環境変数や設定ファイル(`.env`) を利用し、機密情報(APIキー、DBパスワード等)や環境依存の設定値(DBホスト名、ポート番号等)をコード内から分離する