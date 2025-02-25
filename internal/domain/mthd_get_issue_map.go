// internal/domain/strct_master_data.go
package domain

// 必要な getter メソッドを追加 (例)
func (md *MasterData) GetIssueMap() map[string]IssueMaster {
	return md.IssueMap
}

// 他の getter も必要に応じて追加
