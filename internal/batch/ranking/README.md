# estore-trade/internal/batch/ranking

このパッケージは、売買代金ランキングの計算、およびそれに関連するデータアクセス（リポジトリ）を提供する

## 概要
`ranking` パッケージは、バッチ処理の中で呼び出され、立花証券のAPIから株価情報・出来高を取得し、売買代金ランキングを計算する。また、計算されたランキングを基に、取引対象とする銘柄のリストを作成する機能も提供する
ランキング情報をデータベースに永続化するためのリポジトリインターフェース(`RankingRepository`)とその実装も、このパッケージに含まれる

## ファイル構成

-   `calculate.go`: ランキング計算と、取引対象銘柄リスト作成のロジック
-   `repository.go`: `RankingRepository` インターフェースの定義
-   `fact_new_ranking_repository.go`: `RankingRepository` を生成するファクトリ関数
-   `strct_ranking_repository.go`: `RankingRepository` インターフェースを実装する構造体
-   `mthd_save_ranking.go`:  `RankingRepository` の実装 (PostgreSQL への保存処理)
-   `strct_market_data_item.go`: 立花証券APIからのレスポンスを一時的に保持するための構造体
-   `util_calculate_ranking.go`: ランキングを計算
-   `util_create_target_issue_list.go`: ランキングから取引対象銘柄のリストを作成
-   `util_get_market_data.go`: 立花証券APIから株価情報と出来高を取得

### `util_calculate_ranking.go`

-   **`CalculateRanking` 関数**:
    -   立花証券APIクライアント (`tachibana.TachibanaClient`) を使って、全銘柄 (または指定された銘柄) の株価と出来高を取得する
    -   各銘柄の売買代金を計算する (株価 * 出来高)
    -   `domain.Ranking` 構造体のスライスとして、計算されたランキングを返す
    -   エラーが発生した場合は、エラーを返す

### `util_create_target_issue_list.go`

-   **`CreateTargetIssueList` 関数**:
    -   `CalculateRanking` 関数から返されたランキング (`[]domain.Ranking`) と、取得する銘柄数 (`limit`) を受け取る
    -   ランキング上位から `limit` 件の銘柄を抽出し、`domain.TargetIssue` のスライスとして返す

### `iface_ranking_repository.go`

-   **`RankingRepository` インターフェース**:
    -   `SaveRanking(ctx context.Context, ranking []domain.Ranking) error`:  ランキングデータを永続化（DBに保存など）する。
    -   (オプション) `GetRanking(ctx context.Context, date time.Time) ([]domain.Ranking, error)`: 指定された日付のランキングを取得する。
    -   (オプション) `GetLatestRanking(ctx context.Context) ([]domain.Ranking, error)`: 最新のランキングを取得する。

### `fact_new_ranking_repository.go`
- **NewRankingRepository**
    - `RankingRepository`インターフェースを実装した構造体を返す
    - `db *sql.DB`を引数にとり初期化

### `strct_ranking_repository.go`
 - `RankingRepository`インターフェースを実装する構造体
 - `db *sql.DB`をフィールドにもつ

### `mthd_save_ranking.go`
- **SaveRanking**:
    - contextと`[]domain.Ranking`を引数にとりDBに保存

### `strct_market_data_item.go`

-   **`marketDataItem` 構造体**:
     -   立花証券API (`CLMMfdsGetMarketPrice`) から取得したデータを一時的に保持するための構造体。
    -   以下のフィールドを持つ:
        -   `IssueCode`: 銘柄コード
        -   `Price`: 株価
        -   `Volume`: 出来高

### `util_calculate_ranking.go`
-   **`CalculateRanking` 関数**:
    -   contextとtachibanaClientを受け取り、ランキングを計算する

### `util_create_target_issue_list.go`
-   **`CreateTargetIssueList` 関数**:
    -   ランキング情報と、取得する銘柄数を引数で受け取り、銘柄リストを作成

### `util_get_market_data.go`
-   **`getMarketData` 関数**:
    -   contextとtachibanaClient、銘柄コードを受け取り、立花APIにリクエストを送信し、株価と出来高を取得

## 依存関係

-   `context`: コンテキスト
-   `database/sql`: Go の標準 SQL パッケージ
-   `estore-trade/internal/domain`: ドメインモデル (`Ranking`, `TargetIssue`)
-   `estore-trade/internal/infrastructure/persistence/tachibana`: 立花証券 API クライアント

## 特記事項
-  `getMarketData` 関数は、実際には立花証券APIの `CLMMfdsGetMarketPrice`を呼び出して、株価情報を取得する必要がある
-  現状、ランキングを計算するための関数、`CalculateRanking`は、仮実装であり、TODOコメントがある
-  `CreateTargetIssueList`も同様に仮実装
-   エラーハンドリング（リトライ処理、ログ出力など）は適宜追加する必要がある
-   `RankingRepository` インターフェースは、具体的なデータストア (PostgreSQL, MySQL, Redis, ファイルなど) に依存しないように設計され、データストアの変更が容易になる