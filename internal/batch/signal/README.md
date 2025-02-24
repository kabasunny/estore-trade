# estore-trade/internal/batch/signal

このパッケージは、売買シグナルの生成と、生成されたシグナルのデータアクセス（リポジトリ）を提供する

## 概要

`signal` パッケージは、バッチ処理の一部として呼び出され、取引対象の銘柄リストに基づいて売買シグナルを生成する
生成されたシグナルは、データベースに保存される

## ファイル構成

-   `generate.go`: シグナル生成ロジック
-   `repository.go`: `SignalRepository` インターフェースの定義
-   `fact_new_signal_repository.go`: `SignalRepository` を生成するファクトリ関数
-   `strct_signal_repository.go`: `SignalRepository` インターフェースを実装する構造体
-   `mthd_save_signals.go`:  `SignalRepository` の実装 (PostgreSQL への保存処理)

### `util_generate.go`

-   **`Generate` 関数**:
    -   取引対象の銘柄リスト (`[]domain.TargetIssue`) を受け取る
    -   現状は仮実装であり、実際には以下のいずれかの方法でシグナルを生成
        -   外部の Python プロセスを呼び出す (`os/exec` または HTTP リクエストなど)
        -   `internal/autotrading/auto_algorithm` パッケージの `AutoTradingAlgorithm.GenerateSignal()` メソッドを呼び出す (この場合、`GenerateSignal` メソッドの引数と戻り値は、`domain.OrderEvent` ではなく、銘柄リストや関連情報を受け取るように変更する必要がある)
    -   生成されたシグナル (`[]domain.Signal`) を返す
    -   エラーが発生した場合は、エラーを返す

### `iface_signal_repository.go`

-   **`SignalRepository` インターフェース**:
    -   `SaveSignals(ctx context.Context, signals []domain.Signal) error`: 生成されたシグナルを永続化（DBに保存など）する
    -   (オプション) `GetSignals(ctx context.Context, from time.Time, to time.Time) ([]domain.Signal, error)`: 指定された期間のシグナルを取得する
    -   (オプション) `GetLatestSignals(ctx context.Context) ([]domain.Signal, error)`: 最新のシグナルを取得する

### `fact_new_signal_repository.go`
- **NewSignalRepository**:
    - `SignalRepository`インターフェースを実装した構造体を返す
    - `db *sql.DB`を引数にとり初期化

### `strct_signal_repository.go`
 - `SignalRepository`インターフェースを実装する構造体
 - `db *sql.DB`をフィールドにもつ

### `mthd_save_signals.go`

-   **`SaveSignals` メソッド**:
    -   `SignalRepository` インターフェースの実装 (PostgreSQL を使用)
    -   コンテキスト (`context.Context`) とシグナルのスライス (`[]domain.Signal`) を受け取る
    -   データベース接続 (`*sql.DB`) を使用して、シグナルデータを `signals` テーブルに保存する
        -   テーブルが存在しない場合は作成する
        -   トランザクションを使用して、データの整合性を保証する
    -   エラーが発生した場合は、エラーを返す

## 依存関係

-   `context`: コンテキスト
-   `database/sql`: Go の標準 SQL パッケージ
-   `estore-trade/internal/domain`: ドメインモデル (`Signal`)
-   `go.uber.org/zap`: ロギング (Generate関数)

## 特記事項

-   `Generate` 関数は、現状では仮実装であり、実際には外部のシグナル生成器 (Python スクリプトなど) を呼び出すか、`auto_algorithm.AutoTradingAlgorithm.GenerateSignal` メソッドを修正して呼び出す必要がある
-   エラーハンドリング（リトライ処理、ログ出力など）は適宜追加する必要がある
-   `SignalRepository` インターフェースは、具体的なデータストア (PostgreSQL, MySQL, Redis, ファイルなど) に依存しないように設計されている。これにより、データストアの変更が容易になる
-   `signalRepository` 構造体は、`SignalRepository` インターフェースを実装し、PostgreSQL データベースへの接続とデータアクセスロジックをカプセル化する