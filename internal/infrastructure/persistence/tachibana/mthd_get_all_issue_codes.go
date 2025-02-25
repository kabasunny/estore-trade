// internal/infrastructure/persistence/tachibana/mthd_get_all_issue_codes.go
package tachibana

import (
	"context"
	"fmt" //fmtパッケージ
)

func (tc *TachibanaClientImple) GetAllIssueCodes(ctx context.Context) ([]string, error) {
	// MasterDataRepository を使用して DB から取得するように変更
	// return tc.masterDataRepo.GetAllIssueCodes(ctx) //後で実装
	return nil, fmt.Errorf("not implemented yet") //仮
}
