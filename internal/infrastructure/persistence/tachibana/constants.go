// internal/infrastructure/persistence/tachibana/constants.go
package tachibana

// 定数定義 (APIのsCLMIDなど)
const (
	clmidPlaceOrder            = "CLMKabuNewOrder"
	zyoutoekiKazeiCTokutei     = "1"  // 特定口座
	sizyouCToushou             = "00" // 東証
	baibaiKubunBuy             = "3"
	baibaiKubunSell            = "1"
	conditionSashine           = "0" // 指値
	conditionFushinari         = "6" // 不成
	genkinShinyouKubunGenbutsu = "0" // 現物
	orderExpireDay             = "0" // 当日限り

	clmidDownloadMasterData = "CLMEventDownload"
	clmidLogin              = "CLMAuthLoginRequest"
	clmidOrderListDetail    = "CLMOrderListDetail" //注文詳細
	clmidCancelOrder        = "CLMKabuCancelOrder" //注文キャンセル

	clmFdsGetMarketPrice = "CLMMfdsGetMarketPrice" //時価情報
	clmidLogoutRequest   = "CLMAuthLogoutRequest"  // Logout リクエストの sCLMID

)
