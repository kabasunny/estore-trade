package domain

// News はニュース情報（既存の構造体を拡張）
type News struct {
	Title        string   // ニュースタイトル (p_HDL)
	NewsID       string   // ニュースID (p_ID)
	Symbols      []string // 関連銘柄コードリスト (p_ISL)
	CategoryList []string // カテゴリリスト (p_CGL)
	GenreList    []string // ジャンルリスト (p_GRL)
	NewsDate     string   // ニュース日付 (p_DT)  // 追加
	NewsTime     string   // ニュース時刻 (p_TM)  // 追加
	Body         string   // ニュース本文(p_TX) 追加
}
