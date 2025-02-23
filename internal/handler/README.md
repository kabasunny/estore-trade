# estore-trade/internal/handler

このディレクトリは、`estore-trade` アプリケーションにおけるHTTPリクエストのハンドリングを担当

## 概要

`handler` パッケージは、外部からのHTTPリクエストを受け取り、`usecase` パッケージで定義されたビジネスロジックを呼び出し、結果をHTTPレスポンスとして返す
リクエストのバリデーション、認証情報の処理（この例では簡略化）、レスポンスのエンコーディングなども担当

## ファイル構成と詳細説明

-   `fact_new_trading_handler.go`: `TradingHandler` のインスタンスを生成するファクトリ関数
-   `mthd_handle_trade.go`:  `/trade` エンドポイントへのPOSTリクエストを処理する `HandleTrade` メソッドの実装
-   `strct_trading_handler.go`: `TradingHandler` 構造体の定義
-   `util_validate_order_request.go`: 注文リクエストのバリデーションを行う `validateOrderRequest` 関数の定義

### `fact_new_trading_handler.go` (ファクトリ関数)

-   **役割**: `TradingHandler` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewTradingHandler` 関数は、`TradingHandler` のインスタンスを作成し、必要な依存オブジェクト（`TradingUsecase` と `zap.Logger`）を渡す
    -   依存性の注入を行うことで、`TradingHandler` の実装を具体的な依存関係から分離し、テストや変更を容易にする

### `mthd_handle_trade.go` (取引リクエストハンドラ)

-   **役割**: `/trade` エンドポイントへのPOSTリクエストを処理し、注文を実行
-   **詳細**:
    1.  **リクエストボディのデコード**:  JSON形式のリクエストボディを `domain.Order` 構造体に変換し、不正なリクエストボディの場合は、`400 Bad Request` エラーを返す
    2.  **リクエストのバリデーション**: `validateOrderRequest` 関数を使用して、注文リクエストの内容が正しいか検証し、不正なリクエストの場合は、`400 Bad Request` エラーを返す
    3.  **認証情報の取得**: (通常はリクエストヘッダや認証トークンから認証情報を取得するが、今回の実装では認証情報はusecaseでconfigから取得するので、ここでは不要)
    4.  **ユースケースの実行**: `TradingUsecase` の `PlaceOrder` メソッドを呼び出して、注文処理を実行し、注文処理に失敗した場合は、`500 Internal Server Error` エラーを返す
    5.  **レスポンスの作成**: 注文処理の結果（`domain.Order`）をJSON形式にエンコードし、`201 Created` ステータスコードとともにクライアントに返す。レスポンスのエンコーディングに失敗した場合は、`500 Internal Server Error` エラーを返す
    6.  **ログ出力**: 注文が成功した場合は、ログに注文IDなどの情報を出力する

### `strct_trading_handler.go` (構造体定義)

-   **役割**:  `TradingHandler` 構造体を定義
-   **詳細**:
    -   `TradingHandler` 構造体は、HTTPリクエストの処理に必要な依存関係をフィールドとして保持
        -   `tradingUsecase`: 注文処理のビジネスロジックを担当する `TradingUsecase` インターフェース
        -   `logger`: アプリケーションのログ出力を担当する `zap.Logger`

### `util_validate_order_request.go` (リクエストバリデーション)

-   **役割**: 注文リクエストのバリデーションを行う
-   **詳細**:
    -   `validateOrderRequest` 関数は、`domain.Order` 構造体を受け取り、注文数量が正の数であるかなどをチェック
    -   バリデーションルールに違反する場合は、エラーを返す

## 依存関係

-   `encoding/json`: JSONのエンコード/デコード
-   `net/http`: HTTP関連の処理
-   `estore-trade/internal/domain`: ドメインモデル (`Order`)
-   `estore-trade/internal/usecase`: `TradingUsecase` インターフェース
-   `go.uber.org/zap`: ロギングライブラリ

## 特記事項

-   **エラーハンドリング**: 各ステップでエラーが発生する可能性があり、適切にエラーを処理し、クライアントに適切なHTTPステータスコードとエラーメッセージを返す
-   **ロギング**: `zap` ロガーを使用して、リクエストの処理中に発生したイベントやエラーを記録
-   **HTTPステータスコード**:
    -   `201 Created`: 注文が正常に作成された場合に返す
    -   `400 Bad Request`: リクエストボディが不正な場合、またはリクエストのバリデーションに失敗した場合に返す
    -   `500 Internal Server Error`: 注文処理に失敗した場合、またはレスポンスのエンコーディングに失敗した場合に返す
- **認証**: このコードでは認証処理は簡略化されている。　認証情報はusecaseでconfigから取得するので、ここでは不要としている