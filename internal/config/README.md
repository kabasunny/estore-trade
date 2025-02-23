# estore-trade/internal/config

このディレクトリは、`estore-trade` アプリケーションの設定管理を担当

## 概要

`config` パッケージは、アプリケーションの実行に必要な設定情報を環境変数から読み込み、`Config` 構造体に格納
`.env` ファイルを使用して環境変数を設定する

## ファイル構成と詳細説明

-   `strct_config.go`: `Config` 構造体の定義。アプリケーション全体の設定情報を保持
-   `util_login_config.go`: `LoadConfig` 関数の定義。環境変数から設定情報を読み込み、`Config` 構造体のインスタンスを生成

### `strct_config.go` (構造体定義)

-   **役割**: アプリケーションの設定情報を保持する `Config` 構造体を定義
-   **詳細**:
    -   `Config` 構造体は、以下のフィールドを持つ
        -   `TachibanaAPIKey`: 立花証券APIのAPIキー
        -   `TachibanaAPISecret`: 立花証券APIのAPIシークレット
        -   `TachibanaBaseURL`: 立花証券APIのベースURL
        -   `TachibanaUserID`: 立花証券APIのユーザーID
        -   `TachibanaPassword`: 立花証券APIのパスワード
        -   `DBHost`: データベースのホスト名
        -   `DBPort`: データベースのポート番号
        -   `DBUser`: データベースのユーザー名
        -   `DBPassword`: データベースのパスワード
        -   `DBName`: データベース名
        -   `LogLevel`: ログレベル
        -   `EventRid`: イベントのRID (p_rid)
        -   `EventBoardNo`: イベントのボード番号 (p_board_no)
        -   `EventEvtCmd`: イベントのコマンド (p_evt_cmd)
        -   `HTTPPort`: HTTPサーバーのポート番号

### `util_login_config.go` (設定読み込み関数)

-   **役割**: 環境変数から設定情報を読み込み、`Config` 構造体のインスタンスを生成
-   **詳細**:
    -   `LoadConfig` 関数は、以下の処理を行う
        1.  `.env` ファイルの読み込み: `godotenv.Load` を使用して、指定されたパスの `.env` ファイルから環境変数を読み込み、`.env` ファイルが見つからない場合はエラーを返す
        2.  環境変数の取得: `os.Getenv` を使用して、各設定項目に対応する環境変数を取得
        3.  数値への変換: データベースポート (`DB_PORT`) とHTTPポート (`HTTP_PORT`) は文字列として取得されるため、`strconv.Atoi` を使用して整数に変換
            - HTTPポートは、変換に失敗した場合はデフォルト値8080を使用
        4.  `Config` インスタンスの生成: 取得した設定値を `Config` 構造体のフィールドに設定し、インスタンスを返す

## 依存関係

-   `github.com/joho/godotenv`: `.env` ファイルから環境変数を読み込むためのライブラリ
-   `strconv`: 文字列を整数に変換するための標準ライブラリ

## 特記事項

-   **環境変数の設定**: アプリケーションを実行する前に、必要な環境変数を設定する必要がる。`.env` ファイルを使用するか、直接環境変数を設定
-   **エラーハンドリング**: `LoadConfig` 関数は、`.env` ファイルが見つからない場合や、数値変換に失敗した場合にエラーを返し、呼び出し元で適切にエラーハンドリングを行う必要がある
-   **デフォルト値**: HTTPポート (`HTTP_PORT`) は環境変数が設定されてない場合や変換に失敗した場合は、デフォルト値として 8080 が使用され、その他の環境変数は必須
-   **機密情報**: APIキーやパスワードなどの機密情報は、`.env` ファイルに保存し、ソースコードリポジトリには含めない