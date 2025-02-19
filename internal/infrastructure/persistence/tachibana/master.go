// internal/infrastructure/persistence/tachibana/master.go

package tachibana

// SystemStatus システム状態 (必要最低限)
type SystemStatus struct {
	SystemStatusKey string `json:"sSystemStatusKey"` // システム状態ＫＥＹ	001固定
	LoginPermission string `json:"sLoginKyokaKubun"` // ログイン許可区分	0：不許可 1：許可 2：不許可(サービス時間外) 9：管理者のみ可(テスト中)
	SystemState     string `json:"sSystemStatus"`    // システム状態	0：閉局 1：開局 2：一時停止
}

// DateInfo 日付情報 (必要最低限)
type DateInfo struct {
	DateKey           string `json:"sDayKey"`           // 日付ＫＥＹ	001：当日基準 002：翌日基準（夕場）
	PrevBusinessDay1  string `json:"sMaeEigyouDay_1"`   // １営業日前	YYYYMMDD
	TheDay            string `json:"sTheDay"`           // 当日日付	YYYYMMDD
	NextBusinessDay1  string `json:"sYokuEigyouDay_1"`  // 翌１営業日	YYYYMMDD
	StockDeliveryDate string `json:"sKabuUkewatasiDay"` // 株式受渡日	YYYYMMDD

	// 他の日付情報は、必要になったら追加
}

// CallPrice 呼値
type CallPrice struct {
	UnitNumber  int     `json:"sYobineTaniNumber,string"` // 呼値の単位番号
	ApplyDate   string  `json:"sTekiyouDay"`              // 適用日
	Price1      float64 `json:"sKizunPrice_1,string"`     // 基準値段1
	Price2      float64 `json:"sKizunPrice_2,string"`     // 基準値段2
	Price3      float64 `json:"sKizunPrice_3,string"`     // 基準値段3
	Price4      float64 `json:"sKizunPrice_4,string"`     // 基準値段4
	Price5      float64 `json:"sKizunPrice_5,string"`     // 基準値段5
	Price6      float64 `json:"sKizunPrice_6,string"`     // 基準値段6
	Price7      float64 `json:"sKizunPrice_7,string"`     // 基準値段7
	Price8      float64 `json:"sKizunPrice_8,string"`     // 基準値段8
	Price9      float64 `json:"sKizunPrice_9,string"`     // 基準値段9
	Price10     float64 `json:"sKizunPrice_10,string"`    // 基準値段10
	Price11     float64 `json:"sKizunPrice_11,string"`    // 基準値段11
	Price12     float64 `json:"sKizunPrice_12,string"`    // 基準値段12
	Price13     float64 `json:"sKizunPrice_13,string"`    // 基準値段13
	Price14     float64 `json:"sKizunPrice_14,string"`    // 基準値段14
	Price15     float64 `json:"sKizunPrice_15,string"`    // 基準値段15
	Price16     float64 `json:"sKizunPrice_16,string"`    // 基準値段16
	Price17     float64 `json:"sKizunPrice_17,string"`    // 基準値段17
	Price18     float64 `json:"sKizunPrice_18,string"`    // 基準値段18
	Price19     float64 `json:"sKizunPrice_19,string"`    // 基準値段19
	Price20     float64 `json:"sKizunPrice_20,string"`    // 基準値段20
	UnitPrice1  float64 `json:"sYobineTanka_1,string"`    // 呼値単価1
	UnitPrice2  float64 `json:"sYobineTanka_2,string"`    // 呼値単価2
	UnitPrice3  float64 `json:"sYobineTanka_3,string"`    // 呼値単価3
	UnitPrice4  float64 `json:"sYobineTanka_4,string"`    // 呼値単価4
	UnitPrice5  float64 `json:"sYobineTanka_5,string"`    // 呼値単価5
	UnitPrice6  float64 `json:"sYobineTanka_6,string"`    // 呼値単価6
	UnitPrice7  float64 `json:"sYobineTanka_7,string"`    // 呼値単価7
	UnitPrice8  float64 `json:"sYobineTanka_8,string"`    // 呼値単価8
	UnitPrice9  float64 `json:"sYobineTanka_9,string"`    // 呼値単価9
	UnitPrice10 float64 `json:"sYobineTanka_10,string"`   // 呼値単価10
	UnitPrice11 float64 `json:"sYobineTanka_11,string"`   // 呼値単価11
	UnitPrice12 float64 `json:"sYobineTanka_12,string"`   // 呼値単価12
	UnitPrice13 float64 `json:"sYobineTanka_13,string"`   // 呼値単価13
	UnitPrice14 float64 `json:"sYobineTanka_14,string"`   // 呼値単価14
	UnitPrice15 float64 `json:"sYobineTanka_15,string"`   // 呼値単価15
	UnitPrice16 float64 `json:"sYobineTanka_16,string"`   // 呼値単価16
	UnitPrice17 float64 `json:"sYobineTanka_17,string"`   // 呼値単価17
	UnitPrice18 float64 `json:"sYobineTanka_18,string"`   // 呼値単価18
	UnitPrice19 float64 `json:"sYobineTanka_19,string"`   // 呼値単価19
	UnitPrice20 float64 `json:"sYobineTanka_20,string"`   // 呼値単価20
	Decimal1    int     `json:"sDecimal_1,string"`        // 小数点以下の桁数1
	Decimal2    int     `json:"sDecimal_2,string"`        // 小数点以下の桁数2
	Decimal3    int     `json:"sDecimal_3,string"`        // 小数点以下の桁数3
	Decimal4    int     `json:"sDecimal_4,string"`        // 小数点以下の桁数4
	Decimal5    int     `json:"sDecimal_5,string"`        // 小数点以下の桁数5
	Decimal6    int     `json:"sDecimal_6,string"`        // 小数点以下の桁数6
	Decimal7    int     `json:"sDecimal_7,string"`        // 小数点以下の桁数7
	Decimal8    int     `json:"sDecimal_8,string"`        // 小数点以下の桁数8
	Decimal9    int     `json:"sDecimal_9,string"`        // 小数点以下の桁数9
	Decimal10   int     `json:"sDecimal_10,string"`       // 小数点以下の桁数10
	Decimal11   int     `json:"sDecimal_11,string"`       // 小数点以下の桁数11
	Decimal12   int     `json:"sDecimal_12,string"`       // 小数点以下の桁数12
	Decimal13   int     `json:"sDecimal_13,string"`       // 小数点以下の桁数13
	Decimal14   int     `json:"sDecimal_14,string"`       // 小数点以下の桁数14
	Decimal15   int     `json:"sDecimal_15,string"`       // 小数点以下の桁数15
	Decimal16   int     `json:"sDecimal_16,string"`       // 小数点以下の桁数16
	Decimal17   int     `json:"sDecimal_17,string"`       // 小数点以下の桁数17
	Decimal18   int     `json:"sDecimal_18,string"`       // 小数点以下の桁数18
	Decimal19   int     `json:"sDecimal_19,string"`       // 小数点以下の桁数19
	Decimal20   int     `json:"sDecimal_20,string"`       // 小数点以下の桁数20

	// 他の基準値段、呼値単価、小数点以下の桁数は、必要になったら追加
}

// IssueMaster 株式銘柄マスタ (必要最低限)
type IssueMaster struct {
	IssueCode   string `json:"sIssueCode"`         // 銘柄コード
	IssueName   string `json:"sIssueName"`         // 銘柄名称
	TradingUnit int    `json:"sBaibaiTani,string"` // 売買単位
	TokuteiF    string `json:"sTokuteiF"`          // 特定口座対象Ｃ

	// 他の情報は、必要になったら追加
}

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

// OperationStatusKabu 運用ステータス（株）(必要最低限)
type OperationStatusKabu struct {
	ListedMarket string `json:"sZyouzyouSizyou"` // 上場市場
	Unit         string `json:"sUnyouUnit"`      // 運用単位
	Status       string `json:"sUnyouStatus"`    // 運用ステータス
}
