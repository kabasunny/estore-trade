package domain

// SystemStatus システム状態 (必要最低限)
type SystemStatus struct {
	SystemStatusKey string `json:"sSystemStatusKey"` // システム状態ＫＥＹ   001固定
	LoginPermission string `json:"sLoginKyokaKubun"` // ログイン許可区分    0：不許可 1：許可 2：不許可(サービス時間外) 9：管理者のみ可(テスト中)
	SystemState     string `json:"sSystemStatus"`    // システム状態  0：閉局 1：開局 2：一時停止
}
