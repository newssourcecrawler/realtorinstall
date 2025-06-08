package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// HTTPClient is your global client, configured in main.go
var HTTPClient *http.Client

// BaseURL is the address of your API server; can be updated at runtime
var BaseURL = "https://localhost:8443"

// internal storage of the JWT
var authToken string

// SetupHTTPClient configures your HTTPClient with TLS settings
func SetupHTTPClient(cfg *tls.Config) {
	HTTPClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: cfg},
		Timeout:   15 * time.Second,
	}
}

// LoginRequest is the JSON payload for /login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is what the server returns on success
type LoginResponse struct {
	Token string `json:"token"`
}

// Login posts to /login, parses the JWT, and returns it (or an error)
func Login(username, password string) (string, error) {
	// Prepare JSON body
	reqBody := LoginRequest{Username: username, Password: password}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request: %w", err)
	}

	// Send request
	resp, err := HTTPClient.Post(BaseURL+"/login", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		var errObj map[string]string
		_ = json.NewDecoder(resp.Body).Decode(&errObj)
		return "", fmt.Errorf("login failed: %s", errObj["error"])
	}

	// Decode the response
	var lr LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		return "", fmt.Errorf("failed to decode login response: %w", err)
	}

	return lr.Token, nil
}

// SetAuthToken stores the JWT and wraps HTTPClient to include it on every request
func SetAuthToken(token string) {
	authToken = token

	// Wrap the existing Transport so we inject the Bearer header
	base := HTTPClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	HTTPClient.Transport = &authTransport{
		base:  base,
		token: authToken,
	}
}

// authTransport is an http.RoundTripper that adds Authorization headers
type authTransport struct {
	base  http.RoundTripper
	token string
}

func (a *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Inject the header on every outgoing request
	req.Header.Set("Authorization", "Bearer "+a.token)
	return a.base.RoundTrip(req)
}
