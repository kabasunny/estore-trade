import os

# ディレクトリ構成を定義
dirs = [
    "./cmd/trader",
    "./internal/config",
    "./internal/domain",
    "./internal/infrastructure/persistence/tachibana",
    "./internal/infrastructure/database/postgres",
    "./internal/infrastructure/logger/zap",
    "./internal/usecase",
    "./internal/handler",
    "./pkg",
]

# ファイル構成を定義
files = [
    "./cmd/trader/main.go",
    "./internal/config/config.go",
    "./internal/domain/model.go",
    "./internal/domain/repository.go",
    "./internal/infrastructure/persistence/tachibana/tachibana.go",
    "./internal/infrastructure/database/postgres/postgres.go",
    "./internal/infrastructure/logger/zap/zap.go",
    "./internal/usecase/trading.go",
    "./internal/usecase/trading_impl.go",
    "./internal/handler/trading.go",
    "./go.mod",
    "./go.sum",
    "./.env",
]

# ディレクトリを作成
for dir_path in dirs:
    os.makedirs(dir_path, exist_ok=True)

# ファイルを作成
for file_path in files:
    with open(file_path, "w") as f:
        pass

print("ディレクトリおよびファイル構成が作成完了")
