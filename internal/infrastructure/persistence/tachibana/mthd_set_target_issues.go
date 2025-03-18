// internal/infrastructure/persistence/tachibana/mthd_set_target_issues.go
package tachibana

import "context"

// SetTargetIssues は、指定された銘柄コードのみを対象とするようにマスタデータをフィルタリングする
func (tc *TachibanaClientImple) SetTargetIssues(ctx context.Context, issueCodes []string) error {

	//まず、ロックを取得
	tc.mu.Lock()
	defer tc.mu.Unlock()

	// issueCodes をセットに変換 (検索を O(1) にするため)
	issueCodeSet := make(map[string]struct{}, len(issueCodes))
	for _, code := range issueCodes {
		issueCodeSet[code] = struct{}{}
	}

	// issueMap のフィルタリングと関連マップの削除
	for issueCode := range tc.masterData.IssueMap { // すべてのキー（銘柄コード）に対してループ
		if _, exists := issueCodeSet[issueCode]; !exists {
			delete(tc.masterData.IssueMap, issueCode)
			delete(tc.masterData.IssueMarketMap, issueCode)
			delete(tc.masterData.IssueMarketRegulationMap, issueCode)
		}
	}

	// targetIssueCodes の更新
	tc.targetIssueCodesMu.Lock() // 排他制御
	defer tc.targetIssueCodesMu.Unlock()
	tc.targetIssueCodes = issueCodes

	return nil
}
