package tachibana

import "strconv"

// getPNo は p_no を取得し、インクリメントする (スレッドセーフ)
func (tc *TachibanaClientImple) getPNo() string {
	tc.pNoMu.Lock()
	defer tc.pNoMu.Unlock()
	tc.pNo++
	return strconv.FormatInt(tc.pNo, 10)
}
