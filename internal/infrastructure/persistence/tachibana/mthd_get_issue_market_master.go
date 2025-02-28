// internal/infrastructure/persistence/tachibana/mthd_get_issue_market_master.go
package tachibana

import "estore-trade/internal/domain"

// GetIssueMarketMaster は銘柄コードと市場コードに対応する市場情報を返します。
func (tc *TachibanaClientImple) GetIssueMarketMaster(issueCode, marketCode string) (domain.IssueMarketMaster, bool) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()

	// ターゲット銘柄リストが設定されている場合は、チェックを行う
	tc.TargetIssueCodesMu.RLock()
	// ターゲット銘柄リストが空でない、かつ、指定された銘柄コードが含まれていない場合のみ false を返す
	if len(tc.TargetIssueCodes) > 0 && !contains(tc.TargetIssueCodes, issueCode) {
		tc.TargetIssueCodesMu.RUnlock()
		return domain.IssueMarketMaster{}, false // ターゲット銘柄でなければエラー
	}
	tc.TargetIssueCodesMu.RUnlock()

	marketMap, ok := tc.IssueMarketMap[issueCode]
	if !ok {
		return domain.IssueMarketMaster{}, false
	}
	issueMarket, ok := marketMap[marketCode]
	return issueMarket, ok
}
