package tachibana_test

import (
	"estore-trade/internal/domain"
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTestClient(t *testing.T) {

	// テスト用のMasterDataを作成（必要に応じて）
	md := &domain.MasterData{} // ダミーのデータ、またはテスト用のデータ

	t.Run("正常系: クライアントが作成できること", func(t *testing.T) {
		client := tachibana.CreateTestClient(t, md)

		// client が nil でないことを確認
		assert.NotNil(t, client)

		// client が *TachibanaClientImple 型であることを確認
		_, ok := interface{}(client).(*tachibana.TachibanaClientImple)
		assert.True(t, ok)

		// フィールドの内容を出力 (人間による確認用)
		tachibana.PrintClientFields(t, client) // tachibana パッケージの関数を呼び出す

	})
}
