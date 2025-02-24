# estore-trade/internal/domain

このディレクトリは、`estore-trade` アプリケーションのドメインモデルとリポジトリインターフェースを定義する

## 概要

`domain` パッケージは、アプリケーションの中核となるビジネスエンティティ (`Account`、`Order`、`Position`、`OrderEvent`、`Ranking`、`Signal`、`TargetIssue` など) と、それらのデータを永続化するためのリポジトリのインターフェース (`AccountRepository`、`OrderRepository`、`RankingRepository`、`SignalRepository`) を定義する
ドメインモデルは、ビジネスルールやビジネスロジックを表現する上で中心的な役割を果たし、リポジトリインターフェースは、データの永続化方法（データベースの種類など）を抽象化し、ドメイン層が特定のデータベース技術に依存しないようにする

## ファイル構成と詳細説明

-   `iface_account_repository.go`: `AccountRepository` インターフェースの定義、取引アカウント (`Account`) に関するデータ永続化操作を抽象化
-   `iface_order_repository.go`: `OrderRepository` インターフェースの定義、注文 (`Order`) に関するデータ永続化操作を抽象化
-   `strct_account.go`: `Account` 構造体の定義、取引アカウントの情報を表現
-   `strct_order.go`: `Order` 構造体の定義、注文情報を表現
-   `strct_order_event.go`: `OrderEvent` 構造体の定義、注文イベント（注文約定通知、システムステータスなど）の情報を表現
-   `strct_position.go`: `Position` 構造体の定義、保有株の情報を表現
-   `strct_ranking.go`: `Ranking` 構造体の定義、売買代金ランキングの情報を表現
-   `strct_signal.go`: `Signal` 構造体の定義、売買シグナルの情報を表現
-   `strct_target_issue.go`: `TargetIssue` 構造体の定義、取引対象とする銘柄の情報を表現

### `iface_account_repository.go` (アカウントリポジトリインターフェース)

-   **役割**: `Account` データに関する永続化操作を抽象化
-   **詳細**:
    -   `AccountRepository` インターフェースは、以下のメソッドを定義
        -   `GetAccount(ctx context.Context, id string) (*Account, error)`: 指定されたIDのアカウントを取得
        -   `UpdateAccount(ctx context.Context, account *Account) error`: 指定されたアカウントのデータを更新

### `iface_order_repository.go` (注文リポジトリインターフェース)

-   **役割**: `Order` データに関する永続化操作を抽象化
-   **詳細**:
    -   `OrderRepository` インターフェースは、以下のメソッドを定義する
        -   `CreateOrder(ctx context.Context, order *Order) error`: 新しい注文を作成
        -   `GetOrder(ctx context.Context, id string) (*Order, error)`: 指定されたIDの注文を取得
        -   `UpdateOrder(ctx context.Context, order *Order) error`: 指定された注文のデータを更新
        -   `UpdateOrderStatus(ctx context.Context, orderID string, status string) error`: 指定された注文IDの注文のステータスを更新

### `strct_account.go` (アカウント構造体)

-   **役割**: 取引アカウントの情報を表現
-   **詳細**:
    -   `Account` 構造体は、以下のフィールドを持つ
        -   `ID`: アカウントID（ユニーク識別子）
        -   `Balance`: アカウントの現在の残高
        -   `Positions`: ポジションのリスト（取引中または保有中のポジション）
        -   `CreatedAt`: アカウント作成日時
        -   `UpdatedAt`: アカウントの最終更新日時

### `strct_order.go` (注文構造体)

-   **役割**: 注文情報を表現
-   **詳細**:
    -   `Order` 構造体は、以下のフィールドを持つ
        -   `ID`: 注文ID
        -   `Symbol`: 銘柄コード
        -   `Side`: 売買区分 ("buy" または "sell")
        -   `OrderType`: 注文の種類 ("market"、"limit" など)
        -   `Price`: 注文価格
        -   `Quantity`: 注文数量
        -   `Status`: 注文ステータス ("pending"、"filled"、"canceled" など)
        -   `TachibanaOrderID`: 立花証券側の注文ID
        -   `CreatedAt`: 注文作成日時
        -   `UpdatedAt`: 注文最終更新日時

### `strct_order_event.go` (注文イベント構造体)

-   **役割**: 注文イベント（注文約定通知、システムステータスなど）の情報を表現
-   **詳細**:
    -   `OrderEvent` 構造体は、以下のフィールドを持つ
        -   `EventType`: イベントの種類 ("EC"、"NS"、"SS"、"US" など)
        -   `EventNo`: イベント番号 (p_ENO)
        -   `Order`: 更新された注文情報 (ECの場合)
        -   `Timestamp`: イベント発生時刻 (p_date)
        -   `ErrNo`  : エラー番号 (p_errno) (エラーの場合)
        -   `ErrMsg` : エラーメッセージ (p_err) (エラーの場合)

### `strct_position.go` (ポジション構造体)

-   **役割**: 保有株の情報を表現
-   **詳細**:
    -   `Position` 構造体は、以下のフィールドを持つ
        -   `Symbol`: 銘柄コード
        -   `Quantity`: 保有数量
        -   `Price`: 平均取得単価

### `strct_ranking.go` (ランキング構造体)

-   **役割**: 売買代金ランキングの情報を表現
-   **詳細**:
    -   `Ranking` 構造体は、以下のフィールドを持つ
        -   `Rank`: ランキング順位
        -   `IssueCode`: 銘柄コード
        -   `TradingValue`: 売買代金
        -   `CreatedAt`: ランキング生成日時

### `strct_signal.go` (シグナル構造体)

-   **役割**: 売買シグナルの情報を表現
-   **詳細**:
    -   `Signal` 構造体は、以下のフィールドを持つ
        -   `ID`: シグナルID
        -   `IssueCode`: 銘柄コード
        -   `Side`: 売買区分 ("buy" or "sell")
        -   `Priority`: 優先度 (例: 1, 2, 3, ...)
        -   `CreatedAt`: シグナル生成日時

### `strct_target_issue.go` (取引対象銘柄構造体)

-   **役割**: 取引対象とする銘柄の情報を表現
-   **詳細**:
    -   `TargetIssue` 構造体は、以下のフィールドを持つ
        -   `IssueCode`: 銘柄コード

## 依存関係

-   `context`: コンテキスト
-   `time`: 時刻関連の処理

## 特記事項

-   **インターフェース**: リポジトリをインターフェースとして定義することで、具体的なデータアクセス層の実装（データベースの種類など）をドメイン層から分離し、テストの容易性、保守性、拡張性が向上
-   **ドメインモデル**: ドメインモデルは、ビジネスロジックを表現する上で中心的な役割を果たす、必要に応じて、ドメインモデルにメソッドを追加してビジネスルールを実装
-   **不変性**: ドメインオブジェクトは、可能な限り不変 (immutable) に設計することが推奨される、これにより、意図しない状態変更を防ぎ、並行処理における安全性を高める

## 関連リポジトリインターフェース (定義場所: 各Repositoryの実装と同じディレクトリ)

-   `RankingRepository`: `internal/infrastructure/persistence/ranking`
-   `SignalRepository`: `internal/infrastructure/persistence/signal`