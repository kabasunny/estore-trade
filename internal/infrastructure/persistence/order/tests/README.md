# `internal/infrastructure/persistence/order/tests`

このディレクトリには、`order` パッケージのリポジトリ (`orderRepository`) に関するテストコードが含まれる

## テストの方針

*   `database/sql` のモックライブラリとして [sqlmock](https://github.com/DATA-DOG/go-sqlmock) を使用
*   本番コード (`orderRepository`) は、データベース接続に `*sql.DB` を直接使用
*   テストコードは、`sqlmock` を使用して `*sql.DB` の振る舞いを模倣し、実際のデータベース接続は行わない
*   各テスト関数は、`sqlmock.New()` で新しいモックインスタンスを作成し、テストケースごとに期待される SQL クエリと結果を設定
*   テスト対象のメソッドを呼び出し、結果 (戻り値やエラー) を検証
*   `sqlmock.ExpectationsWereMet()` を呼び出し、モックに設定した期待値がすべて満たされたことを確認
*   `orderRepository` の公開メソッド (`CreateOrder`, `GetOrder`, `GetOrdersBySymbolAndStatus`, `UpdateOrder`, `UpdateOrderStatus`, `CancelOrder`) をテスト
*   各メソッドについて、正常系と異常系 (エラーケース、NotFound など) のテストケースを作成

## ファイル構成

*   `create_order_test.go`: `CreateOrder` メソッドのテスト
*   `get_order_test.go`: `GetOrder` メソッドのテスト
*   `get_orders_by_symbol_and_status_test.go`: `GetOrdersBySymbolAndStatus` メソッドのテスト
*   `update_order_test.go`: `UpdateOrder` メソッドのテスト
*   `update_order_status_test.go`: `UpdateOrderStatus` メソッドのテスト
*   `cancel_order_test.go`: `CancelOrder` メソッドのテスト

各テストファイルには、`MockDB` 構造体 (sqlmock.Sqlmock を埋め込んだもの) が定義されている

## テストの実行方法

```bash
go test ./internal/infrastructure/persistence/order/...

で、order パッケージとそのサブディレクトリ (tests を含む) のすべてのテストが実行される
```

## 依存関係

*   github.com/DATA-DOG/go-sqlmock: database/sql のモックライブラリ
*   github.com/stretchr/testify/assert: アサーション (期待値と実際の結果の比較) を行うためのライブラリ