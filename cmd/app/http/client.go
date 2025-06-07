package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

var Client *http.Client

// SetupClient must be called once before any API calls.
func SetupClient(config *tls.Config) {
	Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
		Timeout: 10 * time.Second,
	}
}
