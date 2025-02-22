# estore-trade/internal/infrastructure/persistence/tachibana

このディレクトリは、`estore-trade` アプリケーションにおける立花証券APIとの通信、および関連するデータ処理を担当する `tachibana` パッケージを格納

## 概要

`tachibana` パッケージは、立花証券APIクライアントの実装を提供し、APIとの認証、リクエストの送信、レスポンスの処理、マスタデータの管理、イベントストリームの処理などを行う
このパッケージは、アプリケーションの他の部分（特に `usecase` パッケージ）から利用され、立花証券とのデータ交換を抽象化する

## ファイル構成と詳細説明

-   `constants.go`: 立花証券APIで使用する定数（`sCLMID` など）を定義
-   `fact_new_tachibana_client.go`: `TachibanaClient` の実装を生成するファクトリ関数
-   `iface_client.go`: `TachibanaClient` インターフェースを定義
-   `mthd_*.go`: `TachibanaClient` インターフェースの各メソッドの実装
-   `strct_*.go`: 立花証券APIとのデータ交換に使用する構造体の定義
-   `util_*.go`: ユーティリティ関数群（リトライ処理、Shift-JISデコード、コンテキスト設定など）

### `iface_client.go` (インターフェース定義)

-   **役割**: 立花証券APIクライアントの抽象化されたインターフェースを定義
-   **詳細**:
    -   `TachibanaClient` インターフェースは、立花証券APIとの対話に必要なメソッド（ログイン、注文、注文状況取得、注文キャンセル、イベントストリーム接続、各種URL取得、マスタデータ取得、呼値チェックなど）を宣言
    -   具象型に依存せず、このインターフェースを介してクライアントの機能を利用することで、依存性の注入 (DI) が容易になり、テスト可能性が向上

### `strct_tachibana_client.go` (構造体定義)

-   **役割**: `TachibanaClient` インターフェースを実装する具体的な構造体を定義
-   **詳細**:
    -   `TachibanaClientImple` 構造体は、APIクライアントの実行に必要な情報（ベースURL、APIキー、シークレット、ロガー、ログイン状態、仮想URL、有効期限、排他制御用ミューテックス、連番管理用ミューテックス、マスタデータなど）をフィールドとして保持

### `fact_new_tachibana_client.go` (ファクトリ関数)

-   **役割**: `TachibanaClientImple` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewTachibanaClient` 関数は、`TachibanaClientImple` のインスタンスを作成し、設定情報 (`config.Config`)、ロガー (`zap.Logger`)、および必要な初期値を設定
    -   依存性の注入を行うことで、`TachibanaClientImple` の実装を具体的な依存関係から分離し、テストや変更を容易にする

### `mthd_*.go` (メソッド実装)

各メソッドは、`TachibanaClient` インターフェースで定義された機能を具体的に実装。主要なメソッドは以下の通り

-   `Login`: 立花証券APIにログインし、各種仮想URLを取得・キャッシュ。
-   `PlaceOrder`: 新規注文を送信
-   `GetOrderStatus`: 注文状況を取得
-   `CancelOrder`: 注文をキャンセル
-   `ConnectEventStream`: イベントストリームに接続し、受信したイベントをチャネルに流す（`event_stream.go` で実装）
-   `GetRequestURL`, `GetMasterURL`, `GetPriceURL`, `GetEventURL`: キャッシュされた各種仮想URLを取得
-   `DownloadMasterData`: マスタデータをダウンロードし、内部に保持
-   `GetSystemStatus`, `GetDateInfo`, `GetCallPrice`, `GetIssueMaster`, `GetIssueMarketMaster`, `GetIssueMarketRegulation`, `GetOperationStatusKabu`: 各種マスタデータを取得
-   `CheckPriceIsValid`: 指定された価格が呼値の範囲内であるかを確認
-   `SetTargetIssues`: 指定された銘柄コードのみを対象とするようにマスタデータをフィルタリング
-   `Start`: イベントストリームへの接続を確立し、メッセージ受信ループを開始
-   `Stop`: メッセージ受信ループを停止
-   `sendEvent`: パースされたイベントをusecase層に送信
-   `parseEvent`: 受信したメッセージをパースしてdomain.OrderEventに変換

### `strct_*.go` (構造体定義)

立花証券APIとのデータ交換に使用する構造体を定義

- `CallPrice`: 呼値情報
- `DateInfo`: 日付情報
- `EventStream`: イベントストリームを処理するための構造体
-   `IssueMarketMaster`: 株式銘柄市場マスタ
-   `IssueMarketRegulation`: 株式銘柄別・市場別規制
-   `IssueMaster`: 株式銘柄マスタ
-   `masterDataManager`: マスタデータを一元管理するための構造体
-   `OperationStatusKabu`: 運用ステータス
-   `SystemStatus`: システム状態

### `util_*.go` (ユーティリティ関数)

-   `contains`: スライス内に特定の要素が含まれるかをチェック
-   `formatSDDate`: `time.Time` を特定フォーマットの文字列に変換
-   `isValidPrice`: 注文価格が呼値単位に従っているかをチェック
-   `login`: ログイン処理の共通ロジック
-   `mapToStruct`: `map[string]interface{}` を構造体にマッピング
-   `processResponse`: APIレスポンスを処理し、適切な構造体にデータを格納
-   `retryDo`: リトライ付きでHTTPリクエストを実行
-   `sendRequest`: HTTPリクエストを送信し、レスポンスをデコード（リトライ処理付き）
-   `withContextAndTimeout`: HTTPリクエストにコンテキストとタイムアウトを設定

## 依存関係

-   `context`: コンテキスト管理
-   `encoding/json`: JSONエンコーディング/デコーディング
-   `fmt`: フォーマット関連
-   `io`: I/O関連
-   `math`: 数学関数
-   `net/http`: HTTPクライアント
-   `net/url`: URL解析
-   `strconv`: 文字列変換
-   `sync`: 排他制御
-   `time`: 時間関連
-   `go.uber.org/zap`: ロギング
-   `golang.org/x/text/encoding/japanese`: Shift-JISエンコーディング
-   `golang.org/x/text/transform`: 文字コード変換
-   `estore-trade/internal/config`: 設定情報
-   `estore-trade/internal/domain`: ドメインモデル (Order, OrderEvent など)

## 特記事項

-   **エラーハンドリング**: 各メソッドはエラーを適切に処理し、呼び出し元に返す。APIとの通信ではエラーが発生する可能性があるため、リトライ処理 (`retryDo`) を実装
-   **ロギング**: `zap` ロガーを使用して、重要な操作やエラーに関する情報を記録
-   **Shift-JIS対応**: 立花証券APIはShift-JISを使用するため、レスポンスのデコード時に文字コード変換を行う
-   **仮想URLのキャッシュ**: ログイン時に取得した仮想URLは、`TachibanaClientImple` 構造体にキャッシュされ、有効期限内であれば再利用される
-   **マスタデータの管理**: ダウンロードしたマスタデータは、`TachibanaClientImple` 構造体内のマップに保持され、効率的なアクセスが可能
-   **イベントストリーム**: `ConnectEventStream` メソッド（`event_stream.go` で実装）は、イベントストリームへの接続を確立し、受信したイベントをチャネルに流し、非同期的なイベント処理が可能になる
-   **リトライ処理**: `sendRequest` 関数内で `retryDo` 関数を呼び出し、HTTPリクエストをリトライ付きで実行
-   **呼値チェック**: `CheckPriceIsValid` メソッドで注文価格が呼値の範囲内にあるかを確認
-   **ターゲット銘柄の設定**: `SetTargetIssues` メソッドで、処理対象の銘柄を絞り込む
-   **Long Polling**: `Start` メソッド内で、HTTP GET リクエストをLong Pollingとして送信し、イベントストリームの接続を維持
-   **メッセージ受信ループ**: `Start` メソッド内で、ゴルーチンでメッセージ受信ループを実行し、受信したメッセージをパース、usecase層に通知
-   **停止処理**: `Stop` メソッドで、メッセージ受信ループを停止し、イベントストリームの接続を閉じる