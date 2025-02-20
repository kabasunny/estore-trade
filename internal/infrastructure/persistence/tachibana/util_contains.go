package tachibana

// contains は、スライスに特定の要素が含まれているかどうかをチェックするヘルパー関数
func contains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
