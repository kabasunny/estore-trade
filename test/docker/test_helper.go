// test/docker/test_helper.go
package docker

import (
	"fmt"
	"log"
	"time"

	"estore-trade/internal/domain" // domain パッケージをインポート

	"github.com/ory/dockertest/v3" // docker パッケージをインポート
	"gorm.io/driver/postgres"      // GORM PostgreSQL ドライバ
	"gorm.io/gorm"
)

// ...

func SetupTestDatabase() (*gorm.DB, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	runOpts := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=testuser",
			"POSTGRES_PASSWORD=testpassword",
			"POSTGRES_DB=testdb",
		},
		//Cmd: []string{}, // コマンドを指定する場合はここに記述
		//PortBindings: map[docker.Port][]docker.PortBinding{ //ポートを指定する場合
		//    "5432/tcp": {{HostIP: "0.0.0.0", HostPort: "5432"}},
		//},
	}

	// コンテナを起動
	resource, err := pool.RunWithOptions(runOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}

	var db *gorm.DB

	// ★★★ ここに time.Sleep を追加 ★★★
	// time.Sleep(30 * time.Second) // pcが遅いので待つ

	// 指数バックオフとリトライ回数を設定
	pool.MaxWait = 120 * time.Second

	// Dockerコンテナの準備ができるまで待機 (最大60秒)
	// 指数バックオフを適用
	if err := pool.Retry(func() error {
		dsn := fmt.Sprintf("host=localhost port=%s user=testuser password=testpassword dbname=testdb sslmode=disable", resource.GetPort("5432/tcp")) //修正
		// GORM を使用して PostgreSQL に接続
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			fmt.Println(dsn) // 接続文字列を出力
			return err       // リトライを継続
		}

		// 接続確認 (Ping)
		sqlDB, err := db.DB() // *sql.DB を取得
		if err != nil {
			return err // リトライを継続
		}
		err = sqlDB.Ping()
		if err != nil {
			return err // リトライを継続
		}
		return nil // 成功
	}); err != nil {
		return nil, nil, fmt.Errorf("could not connect to docker: %w", err)
	}

	// AutoMigrate を実行して、テーブルを作成
	err = db.AutoMigrate(&domain.Order{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	cleanup := func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}

	return db, cleanup, nil
}
