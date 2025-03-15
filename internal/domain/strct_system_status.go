package domain

// SystemStatus システム状態 (既存の構造体を拡張)
type SystemStatus struct {
	SystemStatusKey string `json:"systemStatusKey"` // システム状態ＫＥＹ   001固定  //SS, p_SS
	LoginPermission string `json:"loginKyokaKubun"` // ログイン許可区分    0：不許可 1：許可 2：不許可(サービス時間外) 9：管理者のみ可(テスト中) //SS, p_LK
	SystemState     string `json:"systemStatus"`    // システム状態  0：閉局 1：開局 2：一時停止 //SS, p_SS

	ErrNo          string `json:"err_no"`          // エラー番号 //ST,KP,SS,USの共通項目, p_errno
	ErrMsg         string `json:"err_msg"`         // エラーメッセージ //ST,KP,SS,USの共通項目, p_err
	ErrCode        string `json:"err_code"`        // エラーコード  //使用しない
	LoginStatus    string `json:"login_status"`    // ログイン許可区分 (p_LK)  0：不許可 1：許可 2：不許可(サービス時間外) 9：管理者のみ可(テスト中) //SS,USの共通項目, p_LK //LoginPermissionを優先
	MarketCode     string `json:"market_code"`     // 市場コード //US, p_MC
	UpdateStatus   string `json:"update_status"`   // 運用ステータス //US, p_US
	UpdateCategory string `json:"update_category"` // 運用カテゴリー (p_UC)  //US
	UpdateUnit     string `json:"update_unit"`     // 運用ユニット (p_UU)    //US
	UpdateDate     string `json:"update_date"`     //情報更新時間(p_CT)

	UnderlyingAssetCode string // 原資産コード (p_GSCD) //US, 追加
	ProductType         string // 商品種別 (p_SHSB) //US, 追加
	BusinessDayFlag     string // 営業日区分 (p_EDK) //US, 追加
}
