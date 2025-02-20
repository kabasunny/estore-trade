package tachibana

// GetIssueMaster は銘柄コードに対応する銘柄情報を返します。
func (tc *TachibanaClientImple) GetIssueMaster(issueCode string) (IssueMaster, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.targetIssueCodesMu.RLock()
	if len(tc.targetIssueCodes) > 0 {
		if !contains(tc.targetIssueCodes, issueCode) {
			tc.targetIssueCodesMu.RUnlock()
			return IssueMaster{}, false // ターゲット銘柄でなければエラー
		}
	}
	tc.targetIssueCodesMu.RUnlock()

	issue, ok := tc.issueMap[issueCode]
	return issue, ok
}
