### estore-trade

#### ディレクトリ構造と各ファイルの役割概要

- **cmd/**
  - **目的:** 実行可能なアプリケーション（コマンド）のエントリーポイントを格納
  - **trader/**
    - **目的:** 自動売買システム（トレーダー）のメインアプリケーションを配置
    - **main.go**
      - **目的:** プログラムのエントリーポイント
      - **処理概要:**
        - 設定の読み込み (`config.LoadConfig`)
        - ロガーの初期化 (`zapLogger.NewZapLogger`)
        - データベース接続の確立 (`postgres.NewPostgresDB`)
        - OrderRepository と AccountRepository のインスタンス生成
        - 立花証券APIクライアントの初期化 (`tachibana.NewTachibanaClient`)
        - ユースケース層の初期化 (`usecase.NewTradingUsecase`): APIクライアント、ロガー、リポジトリを注入
        - EventStream の初期化と起動: usecase 層から書き込み専用チャネルを取得して渡す
        - EventStream からのイベントを処理するゴルーチンの起動: usecase 層から読み取り専用チャネルを取得して使用
        - HTTPハンドラの初期化 (`handler.NewTradingHandler`): ユースケースとロガーを注入
        - HTTPサーバーの起動 (APIエンドポイント `/trade` を設定)
        - シグナルハンドリング (Graceful Shutdown)

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
      - **目的:** ビジネスエンティティ（`Order`, `Account`, `Position`, `OrderEvent` など）を定義
      - **処理概要:** 自動売買システムで扱うデータ構造（注文、口座、ポジション、注文イベントなど）を定義
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
      - **order_repository.go**
        - **目的:** `OrderRepository` インターフェースの PostgreSQL を使った実装
        - **処理概要:**
          - `Create`, `Read`, `Update`, `Delete`
      - **account_repository.go**
        - **目的:** `AccountRepository` インターフェースの PostgreSQL を使った実装
        - **処理概要:**
          - `Read`, `Update`
      - **tachibana/**
        - **目的:** 立花証券のAPIクライアントの実装
        - **tachibana_client.go**
          - **目的:** 立花証券APIとのインターフェースを定義(抽象化)
          - **処理概要:** `TachibanaClient` インターフェースを定義し、APIとのやり取りに必要なメソッド（`Login`, `PlaceOrder`, `GetOrderStatus`, `CancelOrder`, `ConnectEventStream`）を宣言
        - **client_core.go**
          - **目的:** `TachibanaClientIntImple` 構造体のコア機能（リクエスト送信、共通処理）
          - **処理概要:**
            - `TachibanaClientIntImple` 構造体: APIクライアントの基本情報（ベースURL, APIキー, シークレットキー, ロガー, リクエストURL, 有効期限, 排他制御用ミューテックス, `p_no` 管理）と、マスターデータ（システムステータス、日付情報、呼値、銘柄情報）を保持。
            - `NewTachibanaClient`: `TachibanaClientIntImple` のインスタンスを生成。
            - `Login`:  APIへのログイン処理。キャッシュされたリクエストURLの有効性を確認し、有効であればそれを返す。無効であれば、`client_login.go` の `login` 関数を呼び出して認証処理を行う。
            - `getPNo`: スレッドセーフに `p_no` を取得・インクリメント。
            - `ConnectEventStream`: `event_stream.go` で実装されるため、ここではエラーを返す。
            - `sendRequest`: HTTPリクエストを送信し、レスポンスを処理する共通関数。リクエストのコンテキストとタイムアウト設定、Shift-JISからUTF-8へのデコードを含む。
        - **client_login.go**
          - **目的:** 立花証券APIのログイン処理に特化
          - **処理概要:**
            - `login`: `client_core.go` の `Login` メソッドから呼び出され、実際のAPI認証処理を行う。リトライ処理、レスポンスのステータスコードチェック、Shift-JISからUTF-8への変換、`p_no` の初期値設定、リクエストURLのキャッシュなどを行う。
        - **client_order.go**
          - **目的:** 立花証券APIの注文関連処理（注文、注文状況取得、注文キャンセル）に特化
          - **処理概要:**
            - `PlaceOrder`: 新しい注文を送信。リトライ処理、`p_no` と `p_sd_date` の設定、レスポンスのステータスコードチェック、Shift-JISからUTF-8への変換を行う。
            - `GetOrderStatus`: 注文状況を取得。リトライ処理、`p_no` と `p_sd_date` の設定、レスポンスのステータスコードチェック、Shift-JISからUTF-8への変換を行う。
            - `CancelOrder`: 注文をキャンセル。リトライ処理、`p_no` と `p_sd_date` の設定、レスポンスのステータスコードチェック、Shift-JISからUTF-8への変換を行う。
        -  **client_master_data.go**
           - **目的:** 立花証券APIのマスターデータ関連処理(ダウンロード、データ保持、取得)
           - **処理概要:**
             -  `masterDataManager` 構造体: マスターデータ(システムステータス、日付情報、呼値マップ、銘柄マップ)を保持
             - `DownloadMasterData`: マスターデータをダウンロード。リクエストの作成、レスポンスの処理、`masterDataManager` へのデータ格納、`TachibanaClientIntImple` へのデータコピーを行う。
             - `mapToStruct`: `map[string]interface{}` を構造体にマッピングする汎用関数。
             - `GetSystemStatus`, `GetDateInfo`, `GetCallPrice`, `GetIssueMaster`: マスターデータのゲッターメソッド。
        - **event_stream.go**
          - **目的:** `EventStream` 構造体を定義
          - **処理概要:**
            - `tachibanaClient`: `TachibanaClient` インターフェース
            - `config`: 設定情報
            - `logger`: ロガー
            - `eventCh`: usecase 層にイベントを通知するためのチャネル（送信専用）
            - `stopCh`: ゴルーチンの停止を指示するためのチャネル
            - `conn`: HTTP コネクション
            - `req`: HTTP リクエスト
          - `NewEventStream` 関数で `EventStream` のインスタンスを作成（APIクライアント、設定、ロガー、書き込み専用チャネルを注入）
          - `Start` メソッドで `EVENT I/F` への接続を確立し、メッセージ受信ループを開始（ゴルーチンで実行）
            - HTTP Long Polling を使用
            - 受信メッセージを `parseEvent` でパースし、`sendEvent` で usecase 層に送信
            - エラーハンドリング（リトライ処理など）
          - `Stop` メソッドでメッセージ受信ループを停止
          - `parseEvent` メソッドで受信メッセージをパースし、`domain.OrderEvent` に変換
          - `sendEvent` メソッドで `eventCh` を通じて usecase 層にイベントを送信
         - **master.go**
           - **目的:** 立花証券APIから取得する各種マスターデータの構造体を定義
           -　**処理概要:**
             - `SystemStatus`: システムステータス
             - `DateInfo`: 日付情報
             - `CallPrice`: 呼値
             - `IssueMaster`: 銘柄マスタ(株式)
             - `OperationStatus`: 運用ステータス(未使用)
        - **constants.go**
          - **目的:** 立花証券APIに関連する定数を定義
          - **処理概要:** APIのエンドポイント識別子(`sCLMID`)、取引関連の定数などを定義。
        - **utils.go**
           -　**目的:** ユーティリティ関数を定義
           - **処理概要:**
             - `formatSDDate`: 日付フォーマット関数
             - `withContextAndTimeout`: HTTPリクエストにコンテキストとタイムアウトを設定
             - `retryDo`: HTTPリクエストのリトライ処理 (未使用)
             - `isValidPrice`: 注文価格が呼値の単位に従っているかチェック

- **usecase/**
    - **目的:** アプリケーションのビジネスロジックの実装
    - **trading.go**
        -   **目的:** 取引に関するビジネスロジックのインターフェースを定義
        -   **処理概要:**  `TradingUsecase` インターフェースを定義。`PlaceOrder`（注文）、`GetOrderStatus`（注文状況取得）、`CancelOrder`（注文キャンセル）、`GetEventChannelReader`（イベント受信用チャネル取得）、`GetEventChannelWriter`（イベント送信用チャネル取得）、`HandleOrderEvent`（イベント処理）などのメソッドを持つ。
    - **trading_impl.go**
        -   **目的:** `TradingUsecase` インターフェースの実装
        -   **処理概要:**
            -   `tradingUsecase` 構造体を定義。
            -   `NewTradingUsecase` 関数で、`tradingUsecase` のインスタンスを生成（`TachibanaClient`、`Logger`、リポジトリを注入）。
            -   `PlaceOrder`、`GetOrderStatus`、`CancelOrder` メソッドで、`TachibanaClient` を使用して立花証券APIを呼び出す。
            -   `GetEventChannelReader`、`GetEventChannelWriter` メソッドで、イベントチャネルを返す。
            -   `HandleOrderEvent` メソッドで、受け取ったイベントに応じた処理（データベースの更新など）を行う。

estore-trade/
├── cmd/ # 実行可能なアプリケーション（コマンド）のエントリーポイントを格納
│   │
│   └── trader/ # 自動売買システム（トレーダー）のメインアプリケーションを配置
│       │
│       └── main.go # アプリケーションのエントリーポイント
│           ├── 設定ファイルの読み込み (config.LoadConfig)
│           ├── ロガーの初期化 (zapLogger.NewZapLogger)
│           ├── データベース接続の確立 (postgres.NewPostgresDB)
│           ├── リポジトリの初期化 (persistence.NewOrderRepository, persistence.NewAccountRepository)
│           ├── 外部APIクライアントの初期化 (tachibana.NewTachibanaClient)
│           ├── ユースケース層の初期化 (usecase.NewTradingUsecase)
│           ├── EventStream の初期化と起動 (tachibana.NewEventStream)
│           ├── EventStream からのイベントを処理するゴルーチンの起動
│           ├── HTTPハンドラの初期化 (handler.NewTradingHandler)
│           └── HTTPサーバーの起動とシグナルハンドリング (Graceful Shutdown)
│
├── internal/ # このアプリケーション内でのみ使用されるコードを格納
│   │
│   ├── config/ # 設定管理
│   │   │
│   │   └── config.go
│   │       ├── 環境変数 (.envファイル) から設定情報を読み込み、Config 構造体に格納
│   │       └── Config 構造体: APIキー、DB接続情報、ログレベル、EventStream用パラメータなど
│   │
│   ├── domain/ # ビジネスロジックの中核となるエンティティ（データ構造）とリポジトリインターフェースを定義
│   │   │
│   │   ├── model.go # ビジネスエンティティ (データ構造)
│   │   │   ├── Order (注文)
│   │   │   ├── Account (口座)
│   │   │   ├── Position (保有ポジション)
│   │   │   └── OrderEvent (注文イベント)
│   │   └── repository.go # データアクセス層の抽象化 (インターフェース定義)
│   │       ├── OrderRepository (注文に関するCRUD操作)
│   │       └── AccountRepository (口座に関する操作)
│   │
│   ├── handler/ # HTTPリクエストハンドリング (APIエンドポイント)
│   │   │
│   │   └── trading.go
│   │       ├── TradingHandler 構造体
│   │       ├── NewTradingHandler: TradingHandler インスタンスの作成 (ユースケースとロガーを注入)
│   │       └── HandleTrade: /trade エンドポイントへのリクエスト処理 (ユースケース呼び出し)
│   │
│   ├── infrastructure/ # 外部システムとの連携（データベース、外部API、ロギングなど）に関する具体的な実装を格納
│   │   │
│   │   ├── database/ # データベースとの接続と操作
│   │   │   │
│   │   │   └── postgres/ # PostgreSQLデータベースとの接続
│   │   │       │
│   │   │       └── postgres.go
│   │   │           ├── NewPostgresDB: データベース接続の確立
│   │   │           ├── Close: データベース接続のクローズ
│   │   │           └── DB: *sql.DB インスタンスの取得
│   │   │
│   │   ├── logger/ # ロギング機能の実装
│   │   │   │
│   │   │   └── zapLogger/ # zap ロギングライブラリを使用したロガーの実装
│   │   │       │
│   │   │       └── zapLogger.go
│   │   │           └── NewZapLogger: 設定に基づいた zap.Logger の初期化
│   │   │
│   │   └── persistence/ # データの永続化 (DB, 外部サービス連携)
│   │       │
│   │       ├── account_repository.go # AccountRepository の PostgreSQL 実装
│   │       ├── order_repository.go   # OrderRepository の PostgreSQL 実装
│   │       └── tachibana/ # 立花証券APIクライアント
│   │           │
│   │           ├── client_core.go # TachibanaClientIntImple 構造体のコア機能
│   │           │   ├── TachibanaClientIntImple 構造体: APIクライアント基本情報、マスターデータ保持
│   │           │   ├── NewTachibanaClient: インスタンス生成
│   │           │   ├── Login: APIログイン、仮想URL取得 (キャッシュ、client_login.go 呼び出し)
│   │           │   ├── getPNo: スレッドセーフな p_no 取得・インクリメント
│   │           │   ├── ConnectEventStream: event_stream.go で実装 (ここではエラー)
│   │           │   └── sendRequest: HTTPリクエスト送信、レスポンス処理 (共通関数)
│   │           │
│   │           ├── client_login.go # 立花証券APIのログイン処理
│   │           │   └── login: API認証、リトライ、レスポンス処理、p_no 初期化、URLキャッシュ
│   │           │
│   │           ├── client_order.go # 立花証券APIの注文関連処理
│   │           │   ├── PlaceOrder: 新規注文 (リトライ、p_no/p_sd_date、レスポンス処理)
│   │           │   ├── GetOrderStatus: 注文状況取得 (リトライ、p_no/p_sd_date、レスポンス処理)
│   │           │   └── CancelOrder: 注文キャンセル (リトライ、p_no/p_sd_date、レスポンス処理)
│   │           │
│   │           ├── client_master_data.go # 立花証券APIのマスターデータ関連処理
│   │           │   ├── masterDataManager 構造体: マスターデータ一時保持
│   │           │   ├── DownloadMasterData: ダウンロード、masterDataManager 格納、TachibanaClientIntImple へコピー
│   │           │   ├── mapToStruct: map から構造体へのマッピング
│   │           │   └── GetSystemStatus, GetDateInfo, GetCallPrice, GetIssueMaster: マスターデータ取得
│   │           │
│   │           ├── constants.go # 立花証券API関連の定数
│   │           │
│   │           ├── event_stream.go # EVENT I/F 関連
│   │           │   ├── EventStream 構造体: APIクライアント、設定、ロガー、送信用チャネル、停止用チャネル、HTTP関連
│   │           │   ├── NewEventStream: インスタンス作成 (APIクライアント, 設定, ロガー, 書き込み専用チャネル注入)
│   │           │   ├── Start: 接続確立、メッセージ受信ループ (ゴルーチン, HTTP Long Polling)
│   │           │   ├── Stop: 受信ループ停止
│   │           │   ├── parseEvent: 受信メッセージパース (domain.OrderEvent へ変換)
│   │           │   └── sendEvent: usecase 層へのイベント送信
│   │           │
│   │           ├── master.go # 立花証券APIから取得するマスターデータの構造体定義
│   │           │    ├── SystemStatus: システムステータス
│   │           │    ├── DateInfo:　日付情報
│   │           │    ├── CallPrice: 呼値
│   │           │    ├── IssueMaster: 銘柄マスタ
│   │           │    └── OperationStatus: 運用ステータス
│   │           │
│   │           ├── tachibana_client.go # TachibanaClient インターフェース (APIクライアント抽象化)
│   │           │
│   │           └── utils.go # ユーティリティ関数
│   │               ├── formatSDDate: 日付フォーマット
│   │               ├── withContextAndTimeout: HTTPリクエストにコンテキストとタイムアウト設定
│   │               └── isValidPrice: 注文価格チェック
│   │
│   └── usecase/ # ビジネスロジック
│       
│       ├── trading.go # インターフェース
│       │   
│       │   └── TradingUsecase
│       │       ├── PlaceOrder: 新規注文
│       │       ├── GetOrderStatus: 注文状況取得
│       │       ├── CancelOrder: 注文キャンセル
│       │       ├── GetEventChannelReader: イベント受信用チャネル取得
│       │       ├── GetEventChannelWriter: イベント送信用チャネル取得
│       │       └── HandleOrderEvent: イベント処理
│       │
│       └── trading_impl.go # 実装
│           ├── tradingUsecase 構造体
│           ├── NewTradingUsecase: インスタンス作成 (APIクライアント, ロガー, リポジトリ注入)
│           ├── PlaceOrder: 新規注文 (API, DB)
│           ├── GetOrderStatus: 注文状況取得 (API)
│           ├── CancelOrder: 注文キャンセル (API)
│           ├── GetEventChannelReader: イベント受信用チャネル
│           ├── GetEventChannelWriter: イベント送信用チャネル
│           └── HandleOrderEvent: イベント処理 (DB)