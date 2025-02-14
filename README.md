# estore-trade

## ディレクトリ構造と各ファイルの役割概要

- **cmd/**
  - **目的:** 実行可能なアプリケーション（コマンド）のエントリーポイントを格納
  - **trader/**
    - **目的:** 自動売買システム（トレーダー）のメインアプリケーションを配置
    - **main.go**
      - **目的:** プログラムのエントリーポイント
      - **処理概要:**
        - 設定ファイルの読み込み (`config.LoadConfig`)
        - ロガーの初期化 (`zapLogger.NewZapLogger`)
        - データベース接続の確立 (`postgres.NewPostgresDB`)
        - 外部APIクライアントの初期化 (`tachibana.NewTachibanaClient`)
        - ユースケース層の初期化 (`usecase.NewTradingUsecase`)
        - HTTPハンドラの初期化 (`handler.NewTradingHandler`)
        - HTTPサーバーの起動 (APIエンドポイント `/trade` を設定)

- **internal/**
  - **目的:** このアプリケーション内でのみ使用されるコードを格納
  - **config/**
    - **目的:** 設定管理
    - **config.go**
      - **目的:** 環境変数から設定情報を読み込み、`Config` 構造体に格納
      - **処理概要:**
        - `.env` ファイル（存在する場合）を読み込む (`godotenv.Load`)
        - 環境変数から必要な設定値（APIキー、DB接続情報、ログレベルなど）を取得し、`Config` 構造体を作成して返す
  - **domain/**
    - **目的:** ビジネスロジックの中核となるエンティティ（データ構造）とリポジトリインターフェースを定義
    - **model.go**
      - **目的:** ビジネスエンティティ（`Order`, `Account`, `Position` など）を定義
      - **処理概要:** 自動売買システムで扱うデータ構造（注文、口座、ポジションなど）を定義
    - **repository.go**
      - **目的:** データアクセス層の抽象化 (インターフェース定義)
      - **処理概要:** `OrderRepository` や `AccountRepository` などのインターフェースを定義し、データの永続化に関する操作（CRUDなど）を抽象化
  - **handler/**
    - **目的:** HTTPリクエストのハンドリング（APIエンドポイントの定義とリクエスト/レスポンスの処理）
    - **trading.go**
      - **目的:** 取引関連のHTTPリクエストを処理
      - **処理概要:**
        - `TradingHandler` 構造体を定義
        - `NewTradingHandler` で `TradingHandler` のインスタンスを作成（ユースケースとロガーを注入）
        - `HandleTrade` メソッドで `/trade` エンドポイントへのリクエストを処理:
          - リクエストボディのデコード
          - リクエストのバリデーション
          - ユースケースの実行 (`tradingUsecase.PlaceOrder`)
          - レスポンスの作成と送信
  - **infrastructure/**
    - **目的:** 外部システムとの連携（データベース、外部API、ロギングなど）に関する具体的な実装を格納
    - **database/**
      - **目的:** データベースとの接続と操作
      - **postgres/**
        - **目的:** PostgreSQLデータベースとの接続
        - **postgres.go**
          - **目的:** PostgreSQLデータベースへの接続と基本的な操作を提供
          - **処理概要:**
            - `NewPostgresDB` 関数でデータベース接続を確立し、`PostgresDB` 構造体を返す
            - `Close` メソッドでデータベース接続を閉じる
            - `DB` メソッドで `*sql.DB` インスタンスを取得
    - **logger/**
      - **目的:** ロギング機能の実装
      - **zapLogger/**
        - **目的:** `zap` ロギングライブラリを使用したロガーの実装
        - **zapLogger.go**
          - **目的:** 設定に基づいた `zap.Logger` の初期化
          - **処理概要:**
            - `NewZapLogger` 関数で設定(`config.Config`)を基に`zap.Logger`を生成
            - ログレベル (debug, info, warn, error, dpanic, panic, fatal) に応じた設定
            - 出力形式(ISO8601)や出力先(標準出力、標準エラー出力)の設定
    - **persistence/**
      - **目的:** データの永続化に関する具体的な実装 (特に外部サービスとの連携)
      - **tachibana/**
        - **目的:** 立花証券のAPIクライアントの実装
        - **tachibana_client.go**
          - **目的:** 立花証券APIとのインターフェースを定義(抽象化)
          - **処理概要:** `TachibanaClient` インターフェースを定義し、APIとのやり取りに必要なメソッド（`Login`, `PlaceOrder`, `GetOrderStatus`, `CancelOrder`）を宣言
        - **tachibana_client_impl.go**
          - **目的:** `TachibanaClient` インターフェースの具体的な実装
          - **処理概要:**
            - `TachibanaClientIntImple` 構造体を定義
            - `NewTachibanaClient` 関数で `TachibanaClientIntImple` のインスタンスを作成（設定とロガーを注入）
            - 各メソッド (`Login`, `PlaceOrder`, `GetOrderStatus`, `CancelOrder`) で、立花証券のAPI仕様に沿ってリクエストを作成、送信、レスポンスを処理
  - **usecase/**
    - **目的:** アプリケーションのビジネスロジックを実装
    - **trading.go**
      - **目的:** 取引関連のユースケースのインターフェースを定義
      - **処理概要:** `TradingUsecase` インターフェースを定義し、取引に関する操作（`PlaceOrder`, `GetOrderStatus`, `CancelOrder` など）を宣言
    - **trading_impl.go**
      - **目的:** `TradingUsecase` インターフェースの具体的な実装
      - **処理概要:**
        - `tradingUsecase` 構造体を定義（`TachibanaClient` とロガーを保持）
        - `NewTradingUsecase` 関数で `tradingUsecase` のインスタンスを作成
        - 各メソッド (`PlaceOrder`, `GetOrderStatus`, `CancelOrder`) で、`TachibanaClient` を使用して立花証券APIと連携し、ビジネスロジックを実行

- **pkg/**
  - **目的:** 外部プロジェクトからインポート可能なライブラリを配置する場所 (このプロジェクトでは未使用)




estore-trade
│
├── cmd
│   └── trader
│       └── main.go
│           ├── 設定ファイルの読み込み (config.LoadConfig)
│           ├── ロガーの初期化 (zapLogger.NewZapLogger)
│           ├── データベース接続の確立 (postgres.NewPostgresDB)
│           ├── 外部APIクライアントの初期化 (tachibana.NewTachibanaClient)
│           ├── ユースケース層の初期化 (usecase.NewTradingUsecase)
│           ├── HTTPハンドラの初期化 (handler.NewTradingHandler)
│           └── HTTPサーバーの起動 (APIエンドポイント /trade を設定)
│
├── internal
│   ├── config
│   │   └── config.go
│   │       ├── 環境変数から設定情報を読み込み、Config 構造体に格納
│   │       ├── .env ファイル（存在する場合）を読み込む (godotenv.Load)
│   │       └── 環境変数から設定値を取得し、Config 構造体を作成して返す
│   │
│   ├── domain
│   │   ├── model.go
│   │   │   ├── ビジネスエンティティ (Order, Account, Position) を定義
│   │   │   └── 自動売買システムで扱うデータ構造を定義
│   │   └── repository.go
│   │       ├── データアクセス層の抽象化 (インターフェース定義)
│   │       └── CRUD操作を抽象化 (OrderRepository, AccountRepository)
│   │
│   ├── handler
│   │   └── trading.go
│   │       ├── TradingHandler 構造体を定義
│   │       ├── NewTradingHandler で TradingHandler のインスタンスを作成
│   │       ├── HandleTrade メソッドで /trade エンドポイントへのリクエストを処理
│   │       └── リクエストボディのデコード、バリデーション、ユースケースの実行、レスポンスの作成と送信
│   │
│   ├── infrastructure
│   │   ├── database
│   │   │   ├── postgres
│   │   │   │   └── postgres.go
│   │   │   │       ├── PostgreSQLデータベースへの接続と基本操作
│   │   │   │       ├── NewPostgresDB 関数でデータベース接続を確立
│   │   │   │       ├── Close メソッドでデータベース接続を閉じる
│   │   │   │       └── DB メソッドで *sql.DB インスタンスを取得
│   │   │
│   │   ├── logger
│   │   │   ├── zapLogger
│   │   │   │   └── zapLogger.go
│   │   │   │       ├── zap ロギングライブラリを使用したロガーの実装
│   │   │   │       ├── NewZapLogger 関数で zap.Logger を生成
│   │   │   │       └── ログレベルや出力形式の設定
│   │   │
│   │   ├── persistence
│   │   │   ├── tachibana
│   │   │   │   ├── tachibana_client.go
│   │   │   │   │   ├── 立花証券APIとのインターフェースを定義
│   │   │   │   │   └── TachibanaClient インターフェースを定義
│   │   │   │   └── tachibana_client_impl.go
│   │   │   │       ├── TachibanaClient インターフェースの具体的実装
│   │   │   │       ├── NewTachibanaClient 関数でインスタンスを作成
│   │   │   │       └── 各メソッドで立花証券APIと連携
│   │
│   ├── usecase
│   │   ├── trading.go
│   │   │   ├── TradingUsecase インターフェースを定義
│   │   │   └── 取引に関する操作を宣言
│   │   └── trading_impl.go
│   │       ├── TradingUsecase インターフェースの具体的実装
│   │       ├── tradingUsecase 構造体を定義
│   │       └── 各メソッドでビジネスロジックを実行
│
├── pkg
│   └── 未使用


+-----------------+
| main.go         |
| (cmd/trader)    |
+-----------------+
       | LoadConfig 設定ファイルを読み込み (config.go 参照)
       v
+-----------------+
| config.go       |
| (internal/config) |
+-----------------+
       | NewZapLogger ロガーを初期化 (zapLogger.go 参照)
       v
+-----------------+
| zapLogger.go    |
| (internal/infrastructure/logger/zapLogger) |
+-----------------+
       | NewPostgresDB データベース接続を確立 (postgres.go 参照)
       v
+-----------------+
| postgres.go     |
| (internal/infrastructure/database/postgres) |
+-----------------+
       | NewTachibanaClient APIクライアントを初期化 (tachibana_client.go 参照)
       v
+-----------------+
| tachibana_client.go |
| (internal/infrastructure/persistence/tachibana) |
+-----------------+
       | NewTradingUsecase ユースケース層を初期化 (trading.go 参照)
       v
+-----------------+
| trading.go      |
| (internal/usecase) |
+-----------------+
       | NewTradingHandler HTTPハンドラを初期化し、APIエンドポイント /trade を設定 (trading.go 参照)
       v
+-----------------+
| trading.go      |
| (internal/handler) |
+-----------------+
       | PlaceOrder
       v
+-----------------+
| trading_impl.go |
| (internal/usecase) |
+-----------------+
