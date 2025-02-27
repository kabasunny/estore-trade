# `internal/infrastructure/persistence/master/tests`

このディレクトリには、`master` パッケージのリポジトリ (`masterDataRepository`) に関するテストコードが含まれる

## テストの方針

*   `database/sql` のモックライブラリとして [sqlmock](https://github.com/DATA-DOG/go-sqlmock) を使用
*   本番コード (`masterDataRepository`) は、データベース接続に `*sql.DB` を直接使用
*   テストコードは、`sqlmock` を使用して `*sql.DB` の振る舞いを模倣し、実際のデータベース接続は行わない
*   各テスト関数は、`sqlmock.New()` で新しいモックインスタンスを作成し、テストケースごとに期待される SQL クエリと結果を設定
*   テスト対象のメソッドを呼び出し、結果 (戻り値やエラー) を検証
*   `sqlmock.ExpectationsWereMet()` を呼び出し、モックに設定した期待値がすべて満たされたことを確認
*   `masterDataRepository` の公開メソッド (`SaveMasterData`, `GetAllIssueCodes`, `GetIssueMaster`) をテスト
*   各メソッドについて、正常系と異常系 (エラーケース、NotFound など) のテストケースを作成

## ファイル構成

*   `get_all_issue_codes_test.go`: `GetAllIssueCodes` メソッドのテスト
*   `get_issue_master_test.go`: `GetIssueMaster` メソッドのテスト
*   `save_master_data_test.go`: `SaveMasterData` メソッドのテスト

## テストの実行方法

```bash
go test ./internal/infrastructure/persistence/master/...
```
で、master パッケージとそのサブディレクトリ (tests を含む) のすべてのテストが実行される


## 依存関係
*   github.com/DATA-DOG/go-sqlmock: database/sql のモックライブラリ

*   github.com/stretchr/testify/assert: アサーション (期待値と実際の結果の比較) を行うためのライブラリ