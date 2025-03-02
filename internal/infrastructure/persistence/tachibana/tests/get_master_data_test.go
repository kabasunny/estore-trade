package tachibana_test

import (
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_GetMasterData(t *testing.T) {
	// SetupTestClient を使うと、DownloadMasterData のモックが実行され、
	// masterData が設定される
	client, _ := tachibana.SetupTestClient(t)

	t.Run("MasterData is returned", func(t *testing.T) {
		// SetupTestClient でモックデータが設定されているはず
		// expectedMasterData := &domain.MasterData{}  // 不要な変数を削除

		actualMasterData := client.GetMasterData()

		// ポインタが nil でないことを確認
		assert.NotNil(t, actualMasterData)

		// 必要なフィールドが適切に設定されているか個別に確認 (一部)
		// 例: SystemStatus の検証
		assert.Equal(t, "1", actualMasterData.SystemStatus.SystemState)

		// 例: DateInfo の検証
		assert.Equal(t, "20231101", actualMasterData.DateInfo.TheDay)

		// 例: CallPriceMap の検証 (一部)
		assert.NotEmpty(t, actualMasterData.CallPriceMap)
		assert.Contains(t, actualMasterData.CallPriceMap, "101")

		// 他のフィールドも必要に応じて検証
	})

	t.Run("MasterData is updated after DownloadMasterData", func(t *testing.T) {
		// 新しいモックデータを作成
		newMasterData := &domain.MasterData{
			SystemStatus: domain.SystemStatus{SystemState: "2"}, // 異なる値
			DateInfo:     domain.DateInfo{TheDay: "20231102"},   // 異なる値
			// 他のフィールドも設定...
		}
		// モックデータをクライアントに設定
		tachibana.SetMasterDataForTest(client, newMasterData) // パッケージ関数を呼び出す

		// GetMasterData を呼び出して、新しいデータが返されるか確認
		actualMasterData := client.GetMasterData()
		assert.Equal(t, "2", actualMasterData.SystemStatus.SystemState)
		assert.Equal(t, "20231102", actualMasterData.DateInfo.TheDay)
		// 他のフィールドも同様に確認
	})
	t.Run("Concurrency test", func(t *testing.T) {
		// masterDataの初期値を設定
		initialMasterData := &domain.MasterData{
			SystemStatus: domain.SystemStatus{SystemState: "Initial"},
		}
		tachibana.SetMasterDataForTest(client, initialMasterData) //パッケージ関数を呼び出す

		// ゴルーチン1: masterDataを読み取る
		go func() {
			for i := 0; i < 1000; i++ {
				_ = client.GetMasterData()
			}
		}()

		// ゴルーチン2: masterDataを更新する
		go func() {
			for i := 0; i < 1000; i++ {
				newMasterData := &domain.MasterData{
					SystemStatus: domain.SystemStatus{SystemState: "Updated"},
				}
				tachibana.SetMasterDataForTest(client, newMasterData) //パッケージ関数を呼び出す
			}
		}()
		//少し待つ
		<-time.After(time.Millisecond * 100)
	})
}
