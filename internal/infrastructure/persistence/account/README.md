# estore-trade/internal/infrastructure/persistence/account

このディレクトリは、`estore-trade` アプリケーションにおける口座（アカウント）データの永続化を担当する `account` パッケージを格納

## 概要

`account` パッケージは、`domain.AccountRepository` インターフェースの実装を提供し、データベース (この場合は SQL データベース) との対話を通じて口座データの取得、更新、およびポジション情報の取得を行う
このパッケージは、アプリケーションの他の部分（特に `usecase` パッケージ）から利用され、口座データアクセスの詳細を抽象化する

## ファイル構成と詳細説明

-   `fact_new_account_repository.go`: `accountRepository` 構造体のインスタンスを生成するファクトリ関数
-   `mthd_*.go`: `domain.AccountRepository` インターフェースの各メソッドの実装
-   `strct_account_repository.go`: `domain.AccountRepository` インターフェースを実装する構造体の定義

### `strct_account_repository.go` (構造体定義)

-   **役割**: `domain.AccountRepository` インターフェースを実装する具体的な構造体を定義
-   **詳細**:
    -   `accountRepository` 構造体は、データベース接続 (`*sql.DB`) をフィールドとして保持

### `fact_new_account_repository.go` (ファクトリ関数)

-   **役割**: `accountRepository` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewAccountRepository` 関数は、`accountRepository` のインスタンスを作成し、データベース接続 (`*sql.DB`) を渡す
    -   依存性の注入を行うことで、`accountRepository` の実装を具体的なデータベース接続から分離し、テストや変更を容易にする

### `mthd_*.go` (メソッド実装)

各メソッドは、`domain.AccountRepository` インターフェースで定義された機能を具体的に実装。主要なメソッドは以下の通り

-   `GetAccount`: 指定されたIDの口座情報をデータベースから取得。関連するポジション情報も取得
-   `getPositions`: 指定されたアカウントIDに関連付けられたポジション情報をデータベースから取得 (ヘルパー関数)
-   `UpdateAccount`: 指定された口座情報をデータベースで更新

各メソッドは、SQLクエリを定義し、`database/sql` パッケージの機能を使用してデータベースとの対話を行う。プレースホルダーを使用してSQLインジェクションを防ぎ、エラーハンドリングも適切に行う

## 依存関係

-   `context`: コンテキスト管理
-   `database/sql`: SQLデータベースとの対話
-   `estore-trade/internal/domain`: ドメインモデル (Account, Position) およびリポジトリのインターフェース (AccountRepository)
- `time`: 時間関連

## 特記事項

-   **エラーハンドリング**: 各メソッドはエラーを適切に処理し、呼び出し元に返す
-   **SQLインジェクション対策**: プレースホルダー (`$1`, `$2` など) を使用してSQLクエリを構築
-   **データベースとの対話**: `database/sql` パッケージの機能 (`ExecContext`, `QueryRowContext`, `QueryContext`, `Scan`, `Close`など) を使用して、データベースとの対話を安全かつ効率的に行う
-   **存在しないデータの扱い**: `GetAccount` メソッドでは、指定されたIDの口座が存在しない場合、`sql.ErrNoRows` をチェックし、`nil, nil` を返すことで呼び出し元にデータが存在しないことを伝える
-   **ヘルパー関数の利用**: `GetAccount` メソッド内では、ポジション情報を取得するために `getPositions` ヘルパー関数を呼び出し、コードの重複を避け、可読性と保守性を向上
-   **`time`パッケージの使用**: 口座の更新日時 (`updated_at`) を記録するために、`time.Now()` を使用
- **トランザクション**: このコード例では明示的なトランザクション管理は行われていない、必要に応じて `database/sql` のトランザクション機能 (`BeginTx`, `Commit`, `Rollback`) を使用して、複数のデータベース操作をアトミックに実行する。　特に、`UpdateAccount`とポジションの更新を組み合わせる場合に重要