package auto_model

// SignalとPositionの構造体は仮のものなので、自動売買アルゴリズムに合わせて定義してください。
type Signal struct {
	// シグナルの情報
	Symbol string //例
	Side   string
}
