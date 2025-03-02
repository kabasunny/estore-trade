package tachibana

import "context"

// SetTargetIssues は、指定された銘柄コードのみを対象とするようにマスタデータをフィルタリングする
func (tc *TachibanaClientImple) SetTargetIssues(ctx context.Context, issueCodes []string) error {
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
	tc.targetIssueCodes = issueCodes

	return nil
}

// package tachibana

// import "context"

// // SetTargetIssues は、指定された銘柄コードのみを対象とするようにマスタデータをフィルタリングする
// func (tc *TachibanaClientImple) SetTargetIssues(ctx context.Context, issueCodes []string) error {
// 	tc.mu.Lock()
// 	defer tc.mu.Unlock()

// 	// issueMap のフィルタリング
// 	for issueCode := range tc.issueMap {
// 		if !contains(issueCodes, issueCode) { // ヘルパー関数を使用
// 			delete(tc.issueMap, issueCode)
// 		}
// 	}

// 	// issueMarketMap, issueMarketRegulationMap のフィルタリング (issueMap と同様)
// 	for issueCode := range tc.issueMarketMap {
// 		if !contains(issueCodes, issueCode) {
// 			delete(tc.issueMarketMap, issueCode)
// 			continue // 銘柄コードが削除されたら、その下の市場情報も不要
// 		}
// 		// (特定の市場だけが必要な場合は、ここでさらにフィルタリング)
// 	}

// 	// issueMarketRegulationMap のフィルタリング
// 	for issueCode := range tc.issueMarketRegulationMap {
// 		if !contains(issueCodes, issueCode) {
// 			delete(tc.issueMarketRegulationMap, issueCode)
// 			continue // 銘柄コードが削除されたら、その下の市場情報も不要
// 		}
// 		// (特定の市場だけが必要な場合は、ここでさらにフィルタリング)
// 	}

// 	tc.targetIssueCodesMu.Lock() // 排他制御
// 	tc.targetIssueCodes = issueCodes
// 	tc.targetIssueCodesMu.Unlock()
// 	return nil
// }
