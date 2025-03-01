package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
)

// GetIssueMaster は銘柄コードに対応する銘柄情報を返す
func (tc *TachibanaClientImple) GetIssueMaster(issueCode string) (domain.IssueMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	defer tc.targetIssueCodesMu.RUnlock() // defer を追加して Unlock を確実に行う

	// デバッグ出力 (targetIssueCodes と issueCode の内容を表示)
	// fmt.Printf("DEBUG: targetIssueCodes: %v, issueCode: %s\n", tc.targetIssueCodes, issueCode)

	// issueCodes が空の場合は、無条件に false を返す
	if len(tc.targetIssueCodes) == 0 {
		fmt.Printf("targetIssueCodesが空です")
		return domain.IssueMaster{}, false // issueCodes が空の場合は false
	}

	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			// tc.targetIssueCodesMu.RUnlock()  // defer で Unlock されるので不要
			return domain.IssueMaster{}, false // ターゲット銘柄でなければエラー
		}
	}
	// tc.targetIssueCodesMu.RUnlock() // defer で Unlock されるので不要

	issue, ok := tc.issueMap[issueCode]
	return issue, ok
}
