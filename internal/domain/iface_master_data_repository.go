// internal/domain/iface_master_data_repository.go
package domain

import (
	"context"
)

type MasterDataRepository interface {
	SaveMasterData(ctx context.Context, m *MasterData) error
	GetMasterData(ctx context.Context) (*MasterData, error)
	GetAllIssueCodes(ctx context.Context) ([]string, error)
	GetIssueMaster(ctx context.Context, issueCode string) (*IssueMaster, error) // 変更
	// 他の Get メソッドも必要に応じて追加
}
