package domain

import "time"

// Market は市場情報（既存の構造体を大幅に拡張）
type Market struct {
	Symbol          string    // 銘柄コード (p_colの末尾の情報コードが表す銘柄)
	MarketCode      string    // 市場コード
	RowNumber       string    //行番号
	CurrentPrice    float64   // 現在値 (DPP)  <-- float64 に変更
	PriceStatus     string    // 現在値/前値比較 (DPG)
	Turnover        int       // 売買高 (DV)   <-- int に変更
	BidQuantity     int       // 買気配数量 (BV, GBV1-10) <-- int に変更
	BidPrice        float64   // 買気配値 (QBP, GBP1-10) <-- float64 に変更
	AskQuantity     int       // 売気配数量 (AV, GAV1-10) <-- int に変更
	AskPrice        float64   // 売気配値 (QAP, GAP1-10) <-- float64 に変更
	OpeningPrice    float64   // 始値 (DOP)    <-- float64 に変更
	HighPrice       float64   // 高値 (DHP)    <-- float64 に変更
	LowPrice        float64   // 安値 (DLP)    <-- float64 に変更
	TradingVolume   float64   // 売買代金 (DJ)   <-- float64 に変更
	DailyHighStatus string    // 日通し高値フラグ (DHF)
	DailyLowStatus  string    // 日通し安値フラグ (DLF)
	AskQuoteType    string    //売り気配値種類(QAS)
	BidQuoteType    string    //買い気配値種類(QBS)
	VWAP            float64   //VWAP
	TurnoverRatio   float64   //騰落率(DYRP)
	PreviousClose   float64   //前日終値(PRP)
	PreviousChange  float64   //前日比(DYWP)
	Listing         string    // 所属 (LISS)
	HighTime        time.Time //高値時刻
	LowTime         time.Time //安値時刻
	OpenTime        time.Time //始値時刻
	CurrentTime     time.Time //現在値時刻
	OverSell        int       //売-OVER(QOV)
	UnderBuy        int       //買-UNDER(QUV)
}
