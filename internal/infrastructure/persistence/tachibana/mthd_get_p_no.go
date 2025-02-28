package tachibana

import "strconv"

// getPNo は p_no を取得し、インクリメントする (スレッドセーフ)
func (tc *TachibanaClientImple) getPNo() string {
	tc.PNoMu.Lock()
	defer tc.PNoMu.Unlock()
	tc.PNo++
	return strconv.FormatInt(tc.PNo, 10)
}
