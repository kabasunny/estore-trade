// internal/infrastructure/persistence/tachibana/mthd_get_issue_market_regulation.go
package tachibana

import (
	"estore-trade/internal/domain"
	"fmt"
)

// GetIssueMarketRegulation は銘柄コードと市場コードに対応する規制情報を返す
func (tc *TachibanaClientImple) GetIssueMarketRegulation(issueCode, marketCode string) (domain.IssueMarketRegulation, bool) {
	tc.Mu.RLock()
	defer tc.Mu.RUnlock()

	// まず IssueMarketRegulationMap をチェック
	fmt.Printf("Checking IssueMarketRegulationMap for issueCode=%s\n", issueCode)
	marketMap, ok := tc.IssueMarketRegulationMap[issueCode]
	if !ok {
		fmt.Printf("IssueCode %s not found in IssueMarketRegulationMap\n", issueCode)
		return domain.IssueMarketRegulation{}, false
	}

	fmt.Printf("Found marketMap for %s: %+v\n", issueCode, marketMap)
	issueMarketRegulation, ok := marketMap[marketCode]
	if !ok {
		fmt.Printf("MarketCode %s not found in marketMap for %s\n", marketCode, issueCode)
	}
	// ターゲット銘柄リストによるフィルタリング
	tc.TargetIssueCodesMu.RLock()
	defer tc.TargetIssueCodesMu.RUnlock()
	if len(tc.TargetIssueCodes) > 0 && !contains(tc.TargetIssueCodes, issueCode) {
		return domain.IssueMarketRegulation{}, false
	}

	return issueMarketRegulation, true
}
