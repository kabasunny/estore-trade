# estore-trade/internal/autotrading/auto_algorithm

自動売買アルゴリズムの実装を提供

## 概要

`auto_algorithm` パッケージは、`AutoTradingAlgorithm` 構造体を定義し、シグナル生成、ポジション計算に関するロジックを実装
このパッケージは、純粋なアルゴリズムのロジックに集中しており、外部の依存関係（データベース、APIクライアント、イベントストリームなど）を持たない

## ファイル構成

-   `strct_autotrading_algorithm.go`: `AutoTradingAlgorithm` 構造体の定義
-   `mthd_culculate_position.go`: ポジション計算 (`CalculatePosition`) メソッドの実装
-   `mthd_generate_signal.go`: シグナル生成 (`GenerateSignal`) メソッドの実装

## `strct_autotrading_algorithm.go` (構造体定義)

-   **役割**: 自動売買アルゴリズムの基本構造を定義
-   **詳細**:
    -   `AutoTradingAlgorithm` 構造体は、自動売買アルゴリズムに必要なフィールドやメソッドを定義す（今後アルゴリズム固有のパラメータや状態を保持する）

## `mthd_generate_signal.go` (シグナル生成)

-   **役割**:  `domain.OrderEvent` を受け取り、取引シグナル (`auto_model.Signal`) を生成
-   **詳細**:
    - `GenerateSignal` メソッドは、花証券からのイベント情報（`domain.OrderEvent`）を基に、売買の判断材料となるシグナルを生成
    - シグナル生成の具体的なロジックは、このメソッド内に実装（現在は仮の実装で、常に空のシグナルを返す）

## `mthd_culculate_position.go` (ポジション計算)

-   **役割**: 生成されたシグナル (`auto_model.Signal`) を受け取り、注文するポジション (`auto_model.Position`) を計算
-    **詳細**:
    - `CalculatePosition` メソッドは、`GenerateSignal` メソッドで生成されたシグナルを基に、実際にどれだけの数量を売買するか（ポジションサイズ）を決定
     - ポジション計算には、リスク管理、資金管理のロジックが含まれます（現在は仮の実装で、常に空のポジションを返す）

## 依存関係
- `estore-trade/internal/autotrading/auto_model`: `Signal`、`Position` 構造体を利用
- `estore-trade/internal/domain`: `OrderEvent` 構造体を利用

## 特記事項

-   **純粋なアルゴリズム**: このパッケージは、自動売買の判断ロジックのみに特化
-   **拡張性**: 新しい自動売買アルゴリズムを追加する場合は、`AutoTradingAlgorithm` 構造体を拡張し、`GenerateSignal` や `CalculatePosition` メソッドを実装








autotrading/
├── auto_algorithm/
│   ├── base/                          # 共通の基底クラスやインターフェース
│   │   ├── algorithm.go               # 自動売買アルゴリズムのインターフェース
│   │   ├── position_calculator.go    # ポジション計算のインターフェース/基底クラス
│   │   └── signal_generator.go       # シグナル生成のインターフェース/基底クラス
│   ├── day_trade/                     # デイトレード戦略
│   │   ├── algorithm_day_trade.go     # デイトレードアルゴリズムの具象クラス
│   │   ├── position_day_trade.go      # デイトレード用のポジション計算
│   │   └── signal_day_trade.go        # デイトレード用のシグナル生成
│   ├── swing_trade/                   # スイングトレード戦略
│   │   ├── algorithm_swing_trade.go   # スイングトレードアルゴリズムの具象クラス
│   │   ├── position_swing_trade.go    # スイングトレード用のポジション計算
│   │   └── signal_swing_trade.go      # スイングトレード用のシグナル生成
│   └── factory/                      # アルゴリズムのファクトリ
│        └── algorithm_factory.go        # アルゴリズム生成のファクトリ関数
├── auto_usecase/                    # (既存) 自動売買ユースケース
│   ├── ...
└── strategy_config/                 # (新規) 戦略設定
    └── config.go