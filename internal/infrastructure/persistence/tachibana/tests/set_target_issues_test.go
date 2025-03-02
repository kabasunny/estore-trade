// internal/infrastructure/persistence/tachibana/tests/set_target_issues_test.go
package tachibana_test

import (
	"context"
	. "estore-trade/internal/infrastructure/persistence/tachibana" // package tachibanaをインポート
	"testing"
)

func TestSetTargetIssues(t *testing.T) {
	// SetupTestClient を使ってテスト用のクライアントをセットアップ
	client, _ := SetupTestClient(t)

	// テストケース
	tests := []struct {
		name           string
		targetIssues   []string
		expectedIssues []string // 期待される issueMap のキー
	}{
		{
			name:           "Set target issues to 7974",
			targetIssues:   []string{"7974"},
			expectedIssues: []string{"7974"},
		},
		{
			name:           "Set target issues to 9984",
			targetIssues:   []string{"9984"},
			expectedIssues: []string{"9984"},
		},
		{
			name:           "Set target issues to 7974 and 9984",
			targetIssues:   []string{"7974", "9984"},
			expectedIssues: []string{"7974", "9984"},
		},
		{
			name:           "Set empty target issues",
			targetIssues:   []string{},
			expectedIssues: []string{}, // 空になる
		},
		{
			name:           "Set non-existent issue code",
			targetIssues:   []string{"9999"}, // 存在しない銘柄コード
			expectedIssues: []string{},       // 何も残らない
		},
		{
			name:           "Set a mix of existing and non-existent issue codes",
			targetIssues:   []string{"7974", "9999", "1234"}, // 存在するコードと存在しないコード
			expectedIssues: []string{"7974"},                 // 存在するコードのみ残る
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// SetupTestClientで初期化された状態をリセット
			// clientを作り直す
			client, _ = SetupTestClient(t)

			err := client.SetTargetIssues(ctx, tt.targetIssues)
			if err != nil {
				t.Fatalf("SetTargetIssues returned an error: %v", err)
			}

			// targetIssueCodes の検証
			if len(client.GetMasterData().IssueMap) != len(tt.expectedIssues) {
				t.Errorf("targetIssueCodes length mismatch: expected %d, got %d", len(tt.expectedIssues), len(client.GetMasterData().IssueMap))
			}

			// issueMap の検証
			if len(client.GetMasterData().IssueMap) != len(tt.expectedIssues) {
				t.Errorf("issueMap length mismatch: expected %d, got %d", len(tt.expectedIssues), len(client.GetMasterData().IssueMap))
			}
			for _, expectedCode := range tt.expectedIssues {
				if _, ok := client.GetMasterData().IssueMap[expectedCode]; !ok {
					t.Errorf("issueMap does not contain expected code: %s", expectedCode)
				}
			}

			// issueMarketMap, issueMarketRegulationMap の検証 (issueMap と同様)
			for _, expectedCode := range tt.expectedIssues {
				if _, ok := client.GetMasterData().IssueMarketMap[expectedCode]; !ok {
					t.Errorf("issueMarketMap does not contain expected code: %s", expectedCode)
				}
				if _, ok := client.GetMasterData().IssueMarketRegulationMap[expectedCode]; !ok {
					t.Errorf("issueMarketRegulationMap does not contain expected code: %s", expectedCode)
				}
			}
			for issueCode := range client.GetMasterData().IssueMap {
				if _, ok := client.GetMasterData().IssueMarketMap[issueCode]; !ok {
					t.Errorf("issueMarketMap does not contain issue code that in issueMap: %s", issueCode)
				}
				if _, ok := client.GetMasterData().IssueMarketRegulationMap[issueCode]; !ok {
					t.Errorf("issueMarketRegulationMap does not contain issue code that in issueMap: %s", issueCode)
				}
			}
		})
	}
}
