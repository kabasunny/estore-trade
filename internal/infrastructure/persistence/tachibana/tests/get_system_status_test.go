// internal/infrastructure/persistence/tachibana/get_system_status_test.go
package tachibana_test

import (
	"context"
	"fmt" //fmtパッケージをインポート
	"testing"
	"time" // time パッケージをインポート

	"estore-trade/internal/infrastructure/persistence/tachibana"

	"github.com/stretchr/testify/assert"
)

func TestGetSystemStatus(t *testing.T) {
	client := tachibana.CreateTestClient(t, nil)

	// 異常系のテストケース (DownloadMasterData 前)
	t.Run("異常系: DownloadMasterData前にSystemStatusを取得しようとすると初期値が返ること", func(t *testing.T) {
		// ログイン
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)

		// GetSystemStatus 呼び出し (DownloadMasterData 前)
		systemStatus := client.GetSystemStatus()

		// systemStatus が初期値であることを確認
		assert.Empty(t, systemStatus.SystemStatusKey)
		assert.Empty(t, systemStatus.LoginPermission)
		assert.Empty(t, systemStatus.SystemState)

		// ログアウト (Login したら必ず Logout する)
		err = client.Logout(context.Background())
		assert.NoError(t, err)
	})
	// 1秒待機
	time.Sleep(1 * time.Second)
	// 正常系のテストケース (DownloadMasterData 後)
	t.Run("正常系: DownloadMasterData後にSystemStatusを取得できること", func(t *testing.T) {
		// ログイン
		client := tachibana.CreateTestClient(t, nil) // クライアントを再作成
		err := client.Login(context.Background(), nil)
		assert.NoError(t, err)

		// デバッグプリント追加
		fmt.Printf("masterURL after Login: %s\n", client.GetMasterURLForTest())

		// MasterDataダウンロード
		_, err = client.DownloadMasterData(context.Background())
		assert.NoError(t, err)

		// GetSystemStatus 呼び出し
		systemStatus := client.GetSystemStatus()

		// systemStatus の検証 (具体的な値はデモ環境に依存)
		assert.NotEmpty(t, systemStatus.SystemStatusKey)
		assert.NotEmpty(t, systemStatus.LoginPermission)
		assert.NotEmpty(t, systemStatus.SystemState)

		// ログアウト
		err = client.Logout(context.Background())
		assert.NoError(t, err)
	})
}
