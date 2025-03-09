package tachibana

import (
	"fmt"
	"testing"
	"time"

	"estore-trade/internal/config" // あなたのプロジェクトに合わせて修正
	"estore-trade/internal/domain"
	"path/filepath"
	"runtime"

	"go.uber.org/zap/zaptest"
)

// CreateTestClient はテスト用の TachibanaClientImple インスタンスを作成します。
func CreateTestClient(t *testing.T) *TachibanaClientImple {
	t.Helper()

	// .env ファイルのパスを修正
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get caller information")
	}
	// test_helpers.go から見た .env の相対パス (プロジェクトルート)
	envPath := filepath.Join(filepath.Dir(filename), "../../../../.env") // パスを修正

	// 設定ファイルの読み込み
	cfg, err := config.LoadConfig(envPath) // 絶対パスまたは相対パスを指定
	if err != nil {
		t.Fatalf("Error loading config: %v", err)
	}

	// ロガーの作成 (テスト用)
	logger := zaptest.NewLogger(t) // テストログを出力
	// logger := zap.NewNop()  // ログ出力を抑制する場合

	// テスト用のMasterDataを作成（必要に応じて）
	md := &domain.MasterData{} // ダミーのデータ、またはテスト用のデータ

	// TachibanaClientImple インスタンスの作成
	client := NewTachibanaClient(cfg, logger, md).(*TachibanaClientImple)

	return client
}

// GetClientFields は TachibanaClientImple のフィールドの値を map[string]string 形式で返します。
func GetClientFields(client *TachibanaClientImple) map[string]string {
	return map[string]string{
		"baseURL":    client.baseURL.String(), // *url.URL は String() で文字列に
		"loggined":   fmt.Sprintf("%t", client.loggined),
		"requestURL": client.requestURL,
		"masterURL":  client.masterURL,
		"priceURL":   client.priceURL,
		"eventURL":   client.eventURL,
		"expiry":     client.expiry.Format(time.RFC3339Nano), // time.Time は Format() で文字列に
		"pNo":        fmt.Sprintf("%d", client.pNo),
		// 他のフィールドも必要に応じて追加
	}
}

// PrintClientFields は TachibanaClientImple のフィールドを整形して出力します。
func PrintClientFields(t *testing.T, client *TachibanaClientImple) {
	t.Helper()
	fields := GetClientFields(client) // フィールドの値を取得

	fmt.Println("TachibanaClientImple Fields:")
	for name, value := range fields {
		fmt.Printf("  %s: %s\n", name, value)
	}
}
