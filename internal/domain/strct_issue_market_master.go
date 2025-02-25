package domain

// IssueMarketMaster 株式銘柄市場マスタ
type IssueMarketMaster struct {
	IssueCode               string  `json:"sIssueCode"`             // 銘柄コード
	MarketCode              string  `json:"sZyouzyouSizyou"`        // 上場市場
	PriceRangeMin           float64 `json:"sNehabaMin,string"`      // 値幅下限
	PriceRangeMax           float64 `json:"sNehabaMax,string"`      // 値幅上限
	SinyouC                 string  `json:"sSinyouC"`               // 信用取引区分
	PreviousClose           float64 `json:"sZenzituOwarine,string"` // 前日終値（必要に応じて）
	IssueKubunC             string  `json:"sIssueKubunC"`           // 銘柄区分（必要に応じて）
	ZyouzyouKubun           string  `json:"sZyouzyouKubun"`         // 上場区分 (必要に応じて)
	CallPriceUnitNumber     string  `json:"sYobineTaniNumber"`      // 呼値の単位番号
	CallPriceUnitNumberYoku string  `json:"sYobineTaniNumberYoku"`  // 呼値の単位番号(翌営業日)
	// 他の情報は、必要になったら追加
}
