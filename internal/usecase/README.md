# estore-trade/internal/usecase

このディレクトリは、`estore-trade` アプリケーションにおける取引関連のユースケースを実装

## 概要

`usecase` パッケージは、アプリケーションのビジネスロジックの中核を担い、外部サービス（立花証券API）との連携、データの検証、ドメインモデルの操作、リポジトリを介したデータの永続化などを担当

## ファイル構成と詳細説明

-   `fact_new_trading_usecase.go`: `TradingUsecase` の実装を生成するファクトリ関数
-   `iface_trading.go`: `TradingUsecase` インターフェースを定義
-   `mthd_cancel_order.go`:  注文キャンセル (`CancelOrder`) メソッドの実装
-   `mthd_get_event_channel_reader.go`: イベントチャネルリーダー取得 (`GetEventChannelReader`) メソッドの実装
-   `mthd_get_event_channel_writer.go`: イベントチャネルライター取得 (`GetEventChannelWriter`) メソッドの実装
-   `mthd_get_order_status.go`: 注文状況取得 (`GetOrderStatus`) メソッドの実装
-   `mthd_handle_order_event.go`: 注文イベント処理 (`HandleOrderEvent`) メソッドの実装
-   `mthd_place_order.go`: 注文実行 (`PlaceOrder`) メソッドの実装
-   `strct_trading_usecase.go`: `TradingUsecase` インターフェースの実装となる構造体の定義

### `iface_trading_usecase.go` (インターフェース定義)

-   **役割**: 取引ユースケースの抽象化されたインターフェースを定義
-   **詳細**:
    -   `TradingUsecase` インターフェースは、取引に関する操作（注文、注文状況取得、キャンセル）と、イベント処理に関連するメソッドを宣言
    -   具象型に依存せず、このインターフェースを介してユースケースの機能を利用することで、依存性の注入 (DI) が容易になり、テスト可能性が向上
    -   インターフェースには以下のメソッドが含まれる
        -   `PlaceOrder`: 注文を実行
        -   `GetOrderStatus`: 注文状況を取得
        -   `CancelOrder`: 注文をキャンセル
        -   `GetEventChannelReader`: 注文イベントを受信するためのチャネル（読み取り専用）を取得
        -   `GetEventChannelWriter`: 注文イベントを送信するためのチャネル（書き込み専用）を取得
        -   `HandleOrderEvent`: 注文イベントを処理

### `strct_trading_usecase.go` (構造体定義)

-   **役割**: `TradingUsecase` インターフェースを実装する具体的な構造体を定義
-   **詳細**:
    -   `tradingUsecase` 構造体は、ユースケースの実行に必要な依存関係（外部サービス、リポジトリ、ロガーなど）をフィールドとして保持
        -   `tachibanaClient`: 立花証券APIとの通信を担当するクライアント
        -   `logger`: アプリケーションのログ出力を担当するロガー (zap)
        -   `orderRepo`: 注文情報をデータベースに永続化するためのリポジトリ
        -   `accountRepo`: 口座情報をデータベースに永続化するためのリポジトリ
        -   `eventCh`: 注文イベントを送受信するためのチャネル
        -    `config`: 設定情報

### `fact_new_trading_usecase.go` (ファクトリ関数)

- **役割**: `tradingUsecase` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewTradingUsecase` 関数は、`tradingUsecase` のインスタンスを作成し、必要な依存オブジェクトを渡す
    -   依存性の注入を行うことで、`tradingUsecase` の実装を具体的な依存関係から分離し、テストや変更を容易にする

### `mthd_place_order.go` (注文実行)

-   **役割**: 注文処理のビジネスロジックを実装
-   **詳細**:
    1.  **注文前チェック**:
        -   システムが稼働状態であるかを確認 (立花証券API)
        -   注文する銘柄が有効であるかを確認 (立花証券API)
        -   注文数量が売買単位の倍数であるかを確認 (立花証券API)
        -   注文価格が呼値の範囲内であるかを確認 (立花証券API)
    2.  **注文実行**: 立花証券APIを呼び出して、実際の注文処理を行う
    3.  **注文後処理**:
        -   注文成功後、注文情報をデータベースに保存 (リポジトリ)
        -   データベースへの保存に失敗した場合でも、注文自体は成功しているので、エラーは返さずログに記録

### `mthd_get_order_status.go` (注文状況取得)

- **役割**: 指定された注文IDの注文状況を取得
- **詳細**: 立花証券APIを呼び出して、注文状況を取得し、呼び出し元に返す

### `mthd_cancel_order.go` (注文キャンセル)

-   **役割**: 指定された注文IDの注文をキャンセル
-   **詳細**: 立花証券APIを呼び出して注文をキャンセル

### `mthd_get_event_channel_reader.go` (イベントチャネルリーダー取得)

-   **役割**: 注文イベントを受信するための読み取り専用チャネルを提供
-   **詳細**: `tradingUsecase` が持つイベントチャネル (`eventCh`) の読み取り側を返す

### `mthd_get_event_channel_writer.go` (イベントチャネルライター取得)

-   **役割**: 注文イベントを送信するための書き込み専用チャネルを提供
-   **詳細**: `tradingUsecase` が持つイベントチャネル (`eventCh`) の書き込み側を返す

### `mthd_handle_order_event.go` (注文イベント処理)

-   **役割**: 受信した注文イベントに応じた処理を行う
-   **詳細**:
    -   イベントの種類 (`EventType`) に応じて処理を分岐
        -   `EC` (注文約定通知): 注文情報をデータベース上で更新
        -   `SS` (システムステータス), `US` (運用ステータス), `NS` (ニュース通知): ログ出力
        -   未知のイベントタイプ: 警告ログを出力

## 依存関係

-   `internal/config`: 設定情報
-   `internal/domain`: ドメインモデル (Order, OrderEvent など) およびリポジトリのインターフェース (OrderRepository, AccountRepository)
-   `internal/infrastructure/persistence/tachibana`: 立花証券APIクライアント
-   `go.uber.org/zap`: ロギングライブラリ

## 特記事項

-   **エラーハンドリング**: 各メソッドはエラーを適切に処理し、呼び出し元に返す（特に、外部APIとの連携部分ではエラーが発生する可能性があるため、注意深く処理）
-   **ロギング**: `zap` ロガーを使用して、重要な操作やエラーに関する情報を記録
-   **データベースとの連携**: `orderRepo` を介して注文情報を永続化、API呼び出しが成功した場合は、DB保存の失敗は致命的なエラーとはみなさない（ロギングは行う）
-   **イベント駆動**: 注文イベントはチャネル (`eventCh`) を通じて非同期に処理し、システムの応答性とスケーラビリティが向上
-   **呼値チェック**: `tachibana` パッケージ内の関数を使って、注文価格の妥当性を検証
-   **システム状態**: 注文前に、`tachibanaClient`から取得するシステムステータスが稼働中か確認
-   **売買単位**: 注文数量が、`tachibanaClient`から取得する銘柄情報にある売買単位の倍数か確認