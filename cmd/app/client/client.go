// client/client.go
package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

var HTTPClient *http.Client

func SetupHTTPClient(cfg *tls.Config) {
	HTTPClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: cfg},
		Timeout:   15 * time.Second,
	}
}

func Login(user, pass string) (string, error) { /* POST /login → returns JWT */ }
func SetAuthToken(token string)               { /* set Authorization header on HTTPClient */ }
