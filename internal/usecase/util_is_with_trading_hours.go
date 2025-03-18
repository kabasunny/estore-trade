// internal/usecase/tests/test_util_is_with_trading_hours.go
package usecase

import "time"

func IsWithinTradingHours() bool {
	now := time.Now()
	weekday := now.Weekday()
	hour := now.Hour()
	minute := now.Minute()

	// 平日 (月曜〜金曜) かどうかをチェック
	if weekday < time.Monday || weekday > time.Friday {
		return false
	}

	//午前の取引時間内かチェック
	if hour >= 9 && hour < 11 {
		return true
	}

	//午前の特殊なケースをチェック(11:00-11:30)
	if hour == 11 && minute <= 30 {
		return true
	}

	//午後の取引時間内かチェック
	if hour >= 12 && hour < 15 {
		//12時台は、12:30以降が取引時間内
		if hour == 12 {
			if minute >= 30 {
				return true
			}
		} else { //13時以降
			return true
		}
	}
	return false
}
