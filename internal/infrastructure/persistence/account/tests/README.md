# `internal/infrastructure/persistence/account/tests`

このディレクトリには、`account` パッケージのリポジトリ (`accountRepository`) に関するテストコードが含まれています

## テストの方針

*   `database/sql` のモックライブラリとして [sqlmock](https://github.com/DATA-DOG/go-sqlmock) を使用
*   本番コード (`accountRepository`) は、データベース接続に `*sql.DB` を直接使用
*   テストコードは、`sqlmock` を使用して `*sql.DB` の振る舞いを模倣し、実際のデータベース接続は行わない
*   各テスト関数は、`sqlmock.New()` で新しいモックインスタンスを作成し、テストケースごとに期待される SQL クエリと結果を設定
*   テスト対象のメソッドを呼び出し、結果 (戻り値やエラー) を検証
*   `sqlmock.ExpectationsWereMet()` を呼び出し、モックに設定した期待値がすべて満たされたことを確認
*   `accountRepository` の公開メソッド (`CreateAccount`, `GetAccount`, `GetAccountByUserID`, `UpdateAccount`) をテスト
*   `getPositions` は非公開メソッドだが、`GetAccount`, `GetAccountByUserID` のテストを通じて間接的にテスト
*   各メソッドについて、正常系と異常系 (エラーケース) のテストケースを作成

## ファイル構成

*   `create_account_test.go`: `CreateAccount` メソッドのテスト
*   `get_account_test.go`: `GetAccount` メソッドのテスト
*   `get_by_user_id_test.go`: `GetAccountByUserID` メソッドのテスト
*   `update_account_test.go`: `UpdateAccount` メソッドのテスト

## テストの実行方法

```bash
go test ./internal/infrastructure/persistence/account/...
```
で、account パッケージとそのサブディレクトリ (tests を含む) のすべてのテストが実行される

## 依存関係

*  github.com/DATA-DOG/go-sqlmock: database/sql のモックライブラリ
*  github.com/stretchr/testify/assert: アサーション (期待値と実際の結果の比較) を行うためのライブラリ