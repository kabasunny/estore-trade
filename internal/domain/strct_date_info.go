package domain

// DateInfo 日付情報 (必要最低限)
type DateInfo struct {
	DateKey           string `json:"sDayKey"`           // 日付ＫＥＹ    001：当日基準 002：翌日基準（夕場）
	PrevBusinessDay1  string `json:"sMaeEigyouDay_1"`   // １営業日前    YYYYMMDD
	TheDay            string `json:"sTheDay"`           // 当日日付 YYYYMMDD
	NextBusinessDay1  string `json:"sYokuEigyouDay_1"`  // 翌１営業日    YYYYMMDD
	StockDeliveryDate string `json:"sKabuUkewatasiDay"` // 株式受渡日    YYYYMMDD

	// 他の日付情報は、必要になったら追加
}
