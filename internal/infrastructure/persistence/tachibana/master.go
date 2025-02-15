// internal/infrastructure/persistence/tachibana/master.go

package tachibana

// SystemStatus (仮の定義 - ドキュメントとPythonコードに基づいて修正が必要)
type SystemStatus struct {
	SystemStatusKey string `json:"sSystemStatusKey"`
	LoginPermission string `json:"sLoginKyokaKubun"`
	SystemState     string `json:"sSystemStatus"`
	CreateTime      string `json:"sCreateTime"` // 必要に応じて time.Time に変換
	UpdateTime      string `json:"sUpdateTime"` // 必要に応じて time.Time に変換
	UpdateNumber    string `json:"sUpdateNumber"`
	DeleteFlag      string `json:"sDeleteFlag"`
	DeleteTime      string `json:"sDeleteTime"` // 必要に応じて time.Time に変換
}

// DateInfo (仮の定義 - ドキュメントとPythonコードに基づいて修正が必要)
type DateInfo struct {
	DateKey          string `json:"sDayKey"`
	PrevBusinessDay1 string `json:"sMaeEigyouDay_1"` // 必要に応じて time.Time に変換
	PrevBusinessDay2 string `json:"sMaeEigyouDay_2"` // 必要に応じて time.Time に変換
	// ... 他の日付フィールド ...
	NextBusinessDay10 string `json:"sYokuEigyouDay_10"`
	StockDeliveryDate string `json:"sKabuUkewatasiDay"` // 必要に応じて time.Time に変換
}

// CallPrice (呼値)
type CallPrice struct {
	UnitNumber  int     `json:"sYobineTaniNumber,string"`
	ApplyDate   string  `json:"sTekiyouDay"`
	Price1      float64 `json:"sKizunPrice_1,string"`
	Price2      float64 `json:"sKizunPrice_2,string"`
	Price3      float64 `json:"sKizunPrice_3,string"`
	Price4      float64 `json:"sKizunPrice_4,string"`
	Price5      float64 `json:"sKizunPrice_5,string"`
	Price6      float64 `json:"sKizunPrice_6,string"`
	Price7      float64 `json:"sKizunPrice_7,string"`
	Price8      float64 `json:"sKizunPrice_8,string"`
	Price9      float64 `json:"sKizunPrice_9,string"`
	Price10     float64 `json:"sKizunPrice_10,string"`
	Price11     float64 `json:"sKizunPrice_11,string"`
	Price12     float64 `json:"sKizunPrice_12,string"`
	Price13     float64 `json:"sKizunPrice_13,string"`
	Price14     float64 `json:"sKizunPrice_14,string"`
	Price15     float64 `json:"sKizunPrice_15,string"`
	Price16     float64 `json:"sKizunPrice_16,string"`
	Price17     float64 `json:"sKizunPrice_17,string"`
	Price18     float64 `json:"sKizunPrice_18,string"`
	Price19     float64 `json:"sKizunPrice_19,string"`
	Price20     float64 `json:"sKizunPrice_20,string"`
	UnitPrice1  float64 `json:"sYobineTanka_1,string"`
	UnitPrice2  float64 `json:"sYobineTanka_2,string"`
	UnitPrice3  float64 `json:"sYobineTanka_3,string"`
	UnitPrice4  float64 `json:"sYobineTanka_4,string"`
	UnitPrice5  float64 `json:"sYobineTanka_5,string"`
	UnitPrice6  float64 `json:"sYobineTanka_6,string"`
	UnitPrice7  float64 `json:"sYobineTanka_7,string"`
	UnitPrice8  float64 `json:"sYobineTanka_8,string"`
	UnitPrice9  float64 `json:"sYobineTanka_9,string"`
	UnitPrice10 float64 `json:"sYobineTanka_10,string"`
	UnitPrice11 float64 `json:"sYobineTanka_11,string"`
	UnitPrice12 float64 `json:"sYobineTanka_12,string"`
	UnitPrice13 float64 `json:"sYobineTanka_13,string"`
	UnitPrice14 float64 `json:"sYobineTanka_14,string"`
	UnitPrice15 float64 `json:"sYobineTanka_15,string"`
	UnitPrice16 float64 `json:"sYobineTanka_16,string"`
	UnitPrice17 float64 `json:"sYobineTanka_17,string"`
	UnitPrice18 float64 `json:"sYobineTanka_18,string"`
	UnitPrice19 float64 `json:"sYobineTanka_19,string"`
	UnitPrice20 float64 `json:"sYobineTanka_20,string"`
	Decimal1    int     `json:"sDecimal_1,string"`
	Decimal2    int     `json:"sDecimal_2,string"`
	Decimal3    int     `json:"sDecimal_3,string"`
	Decimal4    int     `json:"sDecimal_4,string"`
	Decimal5    int     `json:"sDecimal_5,string"`
	Decimal6    int     `json:"sDecimal_6,string"`
	Decimal7    int     `json:"sDecimal_7,string"`
	Decimal8    int     `json:"sDecimal_8,string"`
	Decimal9    int     `json:"sDecimal_9,string"`
	Decimal10   int     `json:"sDecimal_10,string"`
	Decimal11   int     `json:"sDecimal_11,string"`
	Decimal12   int     `json:"sDecimal_12,string"`
	Decimal13   int     `json:"sDecimal_13,string"`
	Decimal14   int     `json:"sDecimal_14,string"`
	Decimal15   int     `json:"sDecimal_15,string"`
	Decimal16   int     `json:"sDecimal_16,string"`
	Decimal17   int     `json:"sDecimal_17,string"`
	Decimal18   int     `json:"sDecimal_18,string"`
	Decimal19   int     `json:"sDecimal_19,string"`
	Decimal20   int     `json:"sDecimal_20,string"`
}

//（例）銘柄マスタ
type IssueMaster struct {
	IssueCode      string `json:"sIssueCode"`
	IssueName      string `json:"sIssueName"`
	IssueNameRyaku string `json:"sIssueNameRyaku"`
	IssueNameKana  string `json:"sIssueNameKana"`
	IssueNameEizi  string `json:"sIssueNameEizi"`
	// ... 他のフィールド ...
	MarketCode          string `json:"sSizyouC"`           // 市場コード (例："00")
	TradingUnit         int    `json:"sBaibaiTani,string"` // 売買単位
	CallPriceUnitNumber string `json:"sYobineTaniNumber"`  // 呼値の単位番号
}

//（例）運用ステータス
type OperationStatus struct {
	SystemAccount string `json:"sSystemKouzaKubun"`
	Market        string `json:"sZyouzyouSizyou"` // "00"
	Category      string `json:"sUnyouCategory"`  // "01"
	Unit          string `json:"sUnyouUnit"`      //"0101"
	BusinessDay   string `json:"sEigyouDayC"`     // "0"　など
	Status        string `json:"sUnyouStatus"`
}
