# estore-trade/internal/autotrading/auto_model

自動売買に関連するデータ構造を定義

## 概要

`auto_model` パッケージは、自動売買アルゴリズムで使用される `Signal` や `Position` などのデータ構造、およびそれらに付随するメソッドを定義
これらのデータ構造は、`auto_algorithm` パッケージと `auto_usecase` パッケージの間で情報をやり取りするために使用

## ファイル構成

-   `strct_position.go`: `Position` 構造体の定義
-   `strct_signal.go`: `Signal` 構造体の定義
-   `mthd_should_trade.go`: `Signal` に付随する `ShouldTrade` メソッドの定義

## `strct_signal.go` (シグナル構造体)

-   **役割**: 自動売買の取引シグナルを表現
-   **詳細**:
    -   `Signal` 構造体は、以下のフィールドを持つ（仮の実装）
        -   `Symbol`: 銘柄コード
        -   `Side`: 売買区分 ("buy" or "sell")

## `strct_position.go` (ポジション構造体)

-   **役割**:  計算された注文数量などのポジション情報を表現
-   **詳細**:
    -   `Position` 構造体は、以下のフィールドを持つ（仮の実装）
        -   `Symbol`: 銘柄コード
        -   `Quantity`: 注文数量
        -    `Side`: 売買区分

## `mthd_should_trade.go` (取引判断)

-  **役割**: シグナルに基づいて実際に取引を行うべきかを判断
-   **詳細**:
    - `ShouldTrade` メソッドは、`Signal` 構造体に付随するメソッドであり、シグナルの内容に基づいて取引を実行すべきかどうかを判定（現在は仮の実装で、常に `true` を返す）

## 依存関係

- 依存関係なし

## 特記事項
- SignalとPositionの構造体は仮のものなので、自動売買アルゴリズムに合わせて定義