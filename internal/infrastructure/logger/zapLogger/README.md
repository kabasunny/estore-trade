# `internal/infrastructure/logger/zapLogger`

このディレクトリは、`estore-trade` アプリケーションで使用するロガー (`zap.Logger`) を設定・生成するための `zapLogger` パッケージを提供する

## 概要

`zapLogger` パッケージは、以下の機能を提供する

*   設定 (`config.Config`) に基づいた `zap.Logger` インスタンスの生成 (`NewZapLogger` 関数)
*   ログレベルの設定 (debug, info, warn, error, dpanic, panic, fatal, または空文字列の場合は info)
*   JSON 形式でのログ出力
*   標準出力 (stdout) へのログ出力
*   標準エラー出力 (stderr) へのエラーログ出力
*   ISO8601 形式でのタイムスタンプ出力

## ファイル構成

*   `fact_new_zapLogger.go`: `NewZapLogger` 関数 (ロガーの生成) の実装
*   `fact_new_zapLogger_test.go`: `NewZapLogger` 関数のテスト

## `NewZapLogger` 関数

```go
func NewZapLogger(cfg *config.Config) (*zap.Logger, error)
```

## 引数
*   cfg: *config.Config 型 アプリケーションの設定情報 (ログレベルなど) を保持

## 戻り値
*   *zap.Logger: 設定された zap.Logger のインスタンス
*   error: エラーが発生した場合 (例: 無効なログレベルが指定された場合)

## 処理内容
*   cfg.LogLevel に基づいて、zap.Config を初期化
*   debug: 開発用の設定 (zap.NewDevelopmentConfig())
*   info, warn, error, dpanic, panic, fatal: 本番用の設定 (zap.NewProductionConfig())
*   上記以外 (空文字列を含む): 本番用の設定で、ログレベルは info
*   zap.Config を以下のようにカスタマイズ
*   Encoding: "json" (JSON 形式で出力)
*   EncoderConfig.EncodeTime: zapcore.ISO8601TimeEncoder (ISO8601 形式のタイムスタンプ)
*   OutputPaths: []string{"stdout"} (標準出力にログを出力)
*   ErrorOutputPaths: []string{"stderr"} (標準エラー出力にエラーログを出力)
*   zapCfg.Build() で zap.Logger インスタンスを生成し、返す
*   zapCfg.Build() がエラーを返した場合は、エラーを返す

## 依存関係
*   go.uber.org/zap: 構造化ロギングライブラリ
*   go.uber.org/zap/zapcore: zap のコア機能 (ログレベル、エンコーダなど)
*   estore-trade/internal/config: アプリケーションの設定を保持する構造体
*   go.uber.org/zap/zaptest/observer: zapのテスト用パッケージ

## テスト
*   fact_new_zapLogger_test.go には、NewZapLogger 関数のユニットテストが含まれる
*   zaptest/observer を使用して、ログ出力をキャプチャし、検証
*   正常系 (有効なログレベル) と異常系 (無効なログレベル) のテストケース
*   ログレベル、出力形式 (JSON)、出力先 (標準出力/標準エラー出力) が正しく設定されていることを確認
```
go test ./internal/infrastructure/logger/...
```