package tachibana

// IssueMarketRegulation 株式銘柄別・市場別規制 (必要最低限)
type IssueMarketRegulation struct {
	IssueCode               string `json:"sIssueCode"`               // 銘柄コード
	ListedMarket            string `json:"sZyouzyouSizyou"`          // 上場市場
	StopKubun               string `json:"sTeisiKubun"`              // 停止区分
	GenbutuUrituke          string `json:"sGenbutuUrituke"`          // 現物/売付
	SeidoSinyouSinkiUritate string `json:"sSeidoSinyouSinkiUritate"` // 制度信用/売建
	IppanSinyouSinkiUritate string `json:"sIppanSinyouSinkiUritate"` // 一般信用/売建
	SinyouSyutyuKubun       string `json:"sSinyouSyutyuKubun"`       // 信用一極集中区分（必要に応じて）
	// 他の情報は、必要になったら追加
}
