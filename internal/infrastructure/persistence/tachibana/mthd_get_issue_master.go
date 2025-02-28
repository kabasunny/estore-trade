package tachibana

import "estore-trade/internal/domain"

// GetIssueMaster は銘柄コードに対応する銘柄情報を返します。
func (tc *TachibanaClientImple) GetIssueMaster(issueCode string) (domain.IssueMaster, bool) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.TargetIssueCodesMu.RLock()
	if len(tc.TargetIssueCodes) > 0 {
		if !contains(tc.TargetIssueCodes, issueCode) {
			tc.TargetIssueCodesMu.RUnlock()
			return domain.IssueMaster{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.TargetIssueCodesMu.RUnlock()

	issue, ok := tc.IssueMap[issueCode]
	return issue, ok
}
