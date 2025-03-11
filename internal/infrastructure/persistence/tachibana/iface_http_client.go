// internal/infrastructure/persistence/tachibana/iface_http_client.go
package tachibana

import "net/http"

// HTTPClient interface for mocking http.Client
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
