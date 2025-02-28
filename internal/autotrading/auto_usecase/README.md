# estore-trade/internal/autotrading/auto_usecase

自動売買ユースケースの実装

## 概要

`auto_usecase` パッケージは、`AutoTradingUsecase` インターフェースを定義し、その実装を提供
自動売買の開始・停止、イベント処理、`tradingUsecase` との連携など、自動売買の全体的なフロー制御を担当
実際の注文実行、およびDBの更新は `usecase.TradingUsecase` に委譲

## ファイル構成

-   `fact_new_autotrading_usecase.go`: `AutoTradingUsecase` の実装を生成するファクトリ関数
-   `iface_autotrading.go`: `AutoTradingUsecase` インターフェースの定義
-   `mthd_handle_event.go`: イベント処理 (`HandleEvent`) メソッドの実装
-   `mthd_start.go`: 自動売買開始 (`Start`) メソッドの実装
-   `mthd_stop.go`: 自動売買停止 (`Stop`) メソッドの実装
-   `strct_autotrading_usecase.go`: `AutoTradingUsecase` インターフェースの実装となる構造体の定義

## `fact_new_autotrading_usecase.go` (ファクトリ関数)

-   **役割**: `autoTradingUsecase` 構造体のインスタンスを生成し、依存関係を注入
-   **詳細**:
    -   `NewAutoTradingUsecase` 関数は、`autoTradingUsecase` のインスタンスを作成し、必要な依存オブジェクトを渡す
    -   依存性の注入を行うことで、`autoTradingUsecase` の実装を具体的な依存関係から分離し、テストや変更を容易にする

## `iface_autotrading.go` (インターフェース定義)

-   **役割**: 自動売買ユースケースの抽象化されたインターフェースを定義
-   **詳細**:
    -   `AutoTradingUsecase` インターフェースは、自動売買の開始 (`Start`)、停止 (`Stop`)、イベント処理 (`HandleEvent`) のメソッドを宣言

## `mthd_handle_event.go` (イベント処理)

-   **役割**: `tradingUsecase` から受け取ったイベントを処理し、自動売買アルゴリズム (`auto_algorithm`) を呼び出し、必要に応じて注文 (`tradingUsecase.PlaceOrder`) を実行
-   **詳細**:
    1.  イベントストリームからのイベント (`domain.OrderEvent`) を受信
    2.  `auto_algorithm.GenerateSignal` を呼び出して、シグナルを生成
    3.  `auto_model.Signal.ShouldTrade` で取引を行うか判断
    4.  取引を行う場合、`auto_algorithm.CalculatePosition` を呼び出してポジションを計算
    5.  `tradingUsecase.PlaceOrder` を呼び出して、注文を実行
    6.  エラーが発生した場合は、ログに出力

## `mthd_start.go` (自動売買開始)

-   **役割**: 自動売買を開始
-   **詳細**:
    -   `Start` メソッドは、イベントチャネル (`eventCh`) からのイベント受信を開始するためのゴルーチンを起動

## `mthd_stop.go` (自動売買停止)

-   **役割**: 自動売買を停止
-   **詳細**:
    -    `Stop`メソッドは、自動売買を停止するためのメソッド(現状は空の実装)

## `strct_autotrading_usecase.go` (構造体定義)

-   **役割**: `AutoTradingUsecase` インターフェースを実装する具体的な構造体を定義
-   **詳細**:
    -   `autoTradingUsecase` 構造体は、ユースケースの実行に必要な依存関係をフィールドとして保持
        -   `tradingUsecase`: 注文実行や注文状況取得など、取引に関する機能を提供する `usecase.TradingUsecase` インターフェース
        -   `autoTradingAlgorithm`: シグナル生成やポジション計算を行う `auto_algorithm.AutoTradingAlgorithm` 構造体
        -   `logger`: アプリケーションのログ出力を担当する `zap.Logger`
        -   `config`: アプリケーションの設定情報 (`config.Config`)
        -   `eventCh`: `tradingUsecase` からのイベントを受信するためのチャネル

## 依存関係

-   `estore-trade/internal/config`: 設定情報
-   `estore-trade/internal/domain`: ドメインモデル (`OrderEvent` など)
-   `estore-trade/internal/usecase`: `TradingUsecase` インターフェース
-   `estore-trade/internal/autotrading/auto_algorithm`: `AutoTradingAlgorithm` 構造体
-   `estore-trade/internal/autotrading/auto_model`: `Signal` 構造体
-   `go.uber.org/zap`: ロギングライブラリ

## 特記事項

-   **イベント駆動**: `tradingUsecase` からのイベントをチャネル (`eventCh`) 経由で受信し、非同期に処理
-   **依存性の注入**: 依存オブジェクトは `NewAutoTradingUsecase` ファクトリ関数を通じて注入され、テストや変更を容易にする