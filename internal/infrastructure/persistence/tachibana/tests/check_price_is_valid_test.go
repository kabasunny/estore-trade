// internal/infrastructure/persistence/tachibana/tests/check_price_is_valid_test.go

package tachibana_test

import (
	"estore-trade/internal/infrastructure/persistence/tachibana"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTachibanaClientImple_CheckPriceIsValid(t *testing.T) {
	// SetupTestClient を使うと、DownloadMasterData のモックが実行され、
	// IssueMarketMap, CallPriceMap が設定される
	client, _ := tachibana.SetupTestClient(t)

	// SetupTestClientでモックデータがセットされているので、それを使う
	tests := []struct {
		name      string
		issueCode string
		price     float64
		isNextDay bool
		want      bool
		wantErr   bool
	}{{
		name:      "Valid Price_1 (当日)",
		issueCode: "7974", // 任天堂
		price:     2999,   // 呼値単位 1円
		isNextDay: false,
		want:      true,
		wantErr:   false,
	},
		{
			name:      "Valid Price_2 (当日)",
			issueCode: "7974", // 任天堂
			price:     3005,   // 呼値単位 1円
			isNextDay: false,
			want:      true,
			wantErr:   false,
		}, {
			name:      "Invalid Price_1 (当日)",
			issueCode: "7974",
			price:     3001.5, // 呼値単位に合わない
			isNextDay: false,
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Invalid Price_2 (当日)",
			issueCode: "7974",
			price:     3004, // 呼値単位に合わない
			isNextDay: false,
			want:      false,
			wantErr:   false,
		}, {
			name:      "Valid Price_1 (翌日)",
			issueCode: "7974", // 任天堂 翌日は5円
			price:     111,
			isNextDay: true,
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Valid Price_2 (翌日)",
			issueCode: "7974", // 任天堂 翌日は5円
			price:     10010,
			isNextDay: true,
			want:      true,
			wantErr:   false,
		},
		{
			name:      "Invalid Price_1 (翌日)",
			issueCode: "7974",
			price:     10001, // 呼値単位に合わない
			isNextDay: true,
			want:      false,
			wantErr:   false,
		}, {
			name:      "Invalid Price_2 (翌日)",
			issueCode: "7974",
			price:     10005, // 呼値単位に合わない
			isNextDay: true,
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Issue Not Found",
			issueCode: "9999", // 存在しない銘柄
			price:     100,
			isNextDay: false,
			want:      false,
			wantErr:   true, // エラーを期待
		},
		// 以下、追加のテストケース
		{
			name:      "Price is 0",
			issueCode: "7974",
			price:     0, // 0円は無効
			isNextDay: false,
			want:      false,
			wantErr:   false, // エラーにはならない
		},
		{
			name:      "Price is negative",
			issueCode: "7974",
			price:     -100, // 負の数は無効
			isNextDay: false,
			want:      false,
			wantErr:   false,
		},
		{
			name:      "Price exceeds the max price in call price table",
			issueCode: "7974",
			price:     30000, // 呼値テーブルに存在しない大きな値
			isNextDay: false,
			want:      true, //現状、呼値テーブルの最後の値で割れるのでtrue
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CheckPriceIsValid(tt.issueCode, tt.price, tt.isNextDay)
			if (err != nil) != tt.wantErr {
				t.Errorf("TachibanaClientImple.CheckPriceIsValid() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TachibanaClientImple.CheckPriceIsValid() = %v, want %v", got, tt.want)
			}
			assert.Equal(t, tt.want, got)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
