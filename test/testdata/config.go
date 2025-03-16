// test/testdata/config.go
package testdata

import "estore-trade/internal/config"

func NewTestConfig() *config.Config {
	return &config.Config{
		// テスト用の設定値を設定
		TachibanaBaseURL: "https://example.com",
		// ... 他の設定値 ...
	}
}
