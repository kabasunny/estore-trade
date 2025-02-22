# estore-trade/internal/infrastructure/persistence/order

このディレクトリは、`estore-trade` アプリケーションにおける注文データの永続化を担当する `order` パッケージを格納

## 概要

`order` パッケージは、`domain.OrderRepository` インターフェースの実装を提供し、データベース (この場合は SQL データベース) との対話を通じて注文データの作成、取得、更新を行う
このパッケージは、アプリケーションの他の部分（特に `usecase` パッケージ）から利用され、注文データアクセスの詳細を抽象化

## ファイル構成と詳細説明

-   `fact_new_order_repository.go`: `orderRepository` 構造体のインスタンスを生成するファクトリ関数
-   `mthd_*.go`: `domain.OrderRepository` インターフェースの各メソッドの実装
-   `strct_order_repository.go`: `domain.OrderRepository` インターフェースを実装する構造体の定義

### `strct_order_repository.go` (構造体定義)

-   **役割**: `domain.OrderRepository` インターフェースを実装する具体的な構造体を定義
-   **詳細**:
    -   `orderRepository` 構造体は、データベース接続 (`*sql.DB`) をフィールドとして保持

### `fact_new_order_repository.go` (ファクトリ関数)

-   **役割**: `orderRepository` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewOrderRepository` 関数は、`orderRepository` のインスタンスを作成し、データベース接続 (`*sql.DB`) を渡す
    -   依存性の注入を行うことで、`orderRepository` の実装を具体的なデータベース接続から分離し、テストや変更を容易にする

### `mthd_*.go` (メソッド実装)

各メソッドは、`domain.OrderRepository` インターフェースで定義された機能を具体的に実装。主要なメソッドは以下の通り

-   `CreateOrder`: 新しい注文情報をデータベースに挿入
-   `GetOrder`: 指定されたIDの注文情報をデータベースから取得
-   `UpdateOrder`: 指定された注文情報をデータベースで更新
-  `UpdateOrderStatus`: 指定された注文IDの注文ステータスを更新

各メソッドは、SQLクエリを定義し、`database/sql` パッケージの機能を使用してデータベースとの対話を行う。プレースホルダーを使用してSQLインジェクションを防ぎ、エラーハンドリングも適切に行う

## 依存関係

-   `context`: コンテキスト管理
-   `database/sql`: SQLデータベースとの対話
-   `estore-trade/internal/domain`: ドメインモデル (Order) およびリポジトリのインターフェース (OrderRepository)
-   `time`: 時間関連

## 特記事項

-   **エラーハンドリング**: 各メソッドはエラーを適切に処理し、呼び出し元に返す
-   **SQLインジェクション対策**: プレースホルダー (`$1`, `$2` など) を使用してSQLクエリを構築
-   **データベースとの対話**: `database/sql` パッケージの機能 (`ExecContext`, `QueryRowContext`, `Scan`, `RowsAffected` など) を使用して、データベースとの対話を安全かつ効率的に行う
-   **存在しないデータの扱い**: `GetOrder` メソッドでは、指定されたIDの注文が存在しない場合、`sql.ErrNoRows` をチェックし、`nil, nil` を返すことで呼び出し元にデータが存在しないことを伝える
-   **更新時の存在チェック**: `UpdateOrderStatus`メソッドでは、`RowsAffected`を使って更新対象のデータが存在するか確認し、存在しない場合はエラーを返す
-   **`time`パッケージの使用**: 注文の作成日時 (`created_at`) と更新日時 (`updated_at`) を記録するために、`time.Now()` を使用