// internal/infrastructure/persistence/tachibana/master.go

package tachibana

// SystemStatus (仮の定義 - ドキュメントとPythonコードに基づいて修正が必要)
type SystemStatus struct {
	SystemStatusKey string `json:"sSystemStatusKey"` // システム状態を一意に識別するキー
	LoginPermission string `json:"sLoginKyokaKubun"` // ログイン許可区分 ("0": ログイン不許可, "1": ログイン許可 など)
	SystemState     string `json:"sSystemStatus"`    // システム状態 ("0": サービス停止, "1": サービス中 など)
	CreateTime      string `json:"sCreateTime"`      // 作成日時（必要に応じて time.Time に変換）
	UpdateTime      string `json:"sUpdateTime"`      // 更新日時（必要に応じて time.Time に変換）
	UpdateNumber    string `json:"sUpdateNumber"`    // 更新通番（更新のたびにインクリメントされる番号）
	DeleteFlag      string `json:"sDeleteFlag"`      // 削除フラグ（論理削除用）
	DeleteTime      string `json:"sDeleteTime"`      // 削除時刻（必要に応じて time.Time に変換）
}

// DateInfo (仮の定義 - ドキュメントとPythonコードに基づいて修正が必要)
type DateInfo struct {
	DateKey          string `json:"sDayKey"`         // 日付情報を一意に識別するキー
	PrevBusinessDay1 string `json:"sMaeEigyouDay_1"` // 1営業日前の日付（必要に応じて time.Time に変換）
	PrevBusinessDay2 string `json:"sMaeEigyouDay_2"` // 2営業日前の日付（必要に応じて time.Time に変換）
	// ... 他の日付フィールド ...
	NextBusinessDay10 string `json:"sYokuEigyouDay_10"` // 10営業日後の日付（必要に応じて time.Time に変換）
	StockDeliveryDate string `json:"sKabuUkewatasiDay"` // 株式の受渡日（必要に応じて time.Time に変換）
}

// CallPrice (呼値)
type CallPrice struct {
	UnitNumber  int     `json:"sYobineTaniNumber,string"` // 呼値の単位番号
	ApplyDate   string  `json:"sTekiyouDay"`              // 呼値の適用日
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
}

// IssueMaster (株式の銘柄に関する基本情報)
type IssueMaster struct {
	IssueCode           string `json:"sIssueCode"`         // 銘柄コード (例: "7203" (トヨタ自動車))
	IssueName           string `json:"sIssueName"`         // 銘柄名 (例: "トヨタ自動車")
	IssueNameRyaku      string `json:"sIssueNameRyaku"`    // 銘柄略称
	IssueNameKana       string `json:"sIssueNameKana"`     // 銘柄名（カナ）
	IssueNameEizi       string `json:"sIssueNameEizi"`     // 銘柄名（英字）
	MarketCode          string `json:"sSizyouC"`           // 市場コード (例: "00" (東証))
	TradingUnit         int    `json:"sBaibaiTani,string"` // 売買単位 (例: 100 株)
	CallPriceUnitNumber string `json:"sYobineTaniNumber"`  // 呼値の単位番号 (CallPrice と紐付けに使用)
}

// OperationStatus (市場の運用状態)
type OperationStatus struct {
	SystemAccount string `json:"sSystemKouzaKubun"` // システム口座区分
	Market        string `json:"sZyouzyouSizyou"`   // 上場市場 (例: "00" (東証))
	Category      string `json:"sUnyouCategory"`    // 運用カテゴリ (例: "01" (株式))
	Unit          string `json:"sUnyouUnit"`        // 運用単位
	BusinessDay   string `json:"sEigyouDayC"`       // 営業日区分
	Status        string `json:"sUnyouStatus"`      // 運用ステータス (注文受付中、立会終了 など)
}
