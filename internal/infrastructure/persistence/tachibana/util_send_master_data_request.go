// internal/infrastructure/persistence/tachibana/mthd_download_master_data.go
package tachibana

import (
	"fmt"
	"io"
	"net/http"
)

// sendMasterDataRequest は、マスタデータ取得専用のHTTPリクエスト関数
func sendMasterDataRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned non-200 status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
