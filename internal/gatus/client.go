package gatus

import (
	"net/http"
	"net/url"
	"time"

	"gatus-status/internal/config"

	sdk "github.com/TwiN/gatus-sdk"
)

const timeout = 10 * time.Second

// NewClient creates a Gatus SDK client, adding Basic Auth if configured
func NewClient(cfg config.GatusConfig) *sdk.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 10

	var roundTripper http.RoundTripper = transport
	if cfg.Username != "" || cfg.Password != "" {
		baseURL, _ := url.Parse(cfg.URL)
		roundTripper = &basicAuthTransport{username: cfg.Username, password: cfg.Password, host: baseURL.Host, base: transport}
	}

	return sdk.NewClient(cfg.URL, sdk.WithHTTPClient(&http.Client{Timeout: timeout, Transport: roundTripper}))
}

// basicAuthTransport sets Basic Auth on requests to the Gatus host, as the SDK has no option for it.
// Requests to other hosts (e.g. after a redirect) are sent without credentials.
type basicAuthTransport struct {
	username, password string
	host               string
	base               http.RoundTripper
}

func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == t.host {
		req = req.Clone(req.Context())
		req.SetBasicAuth(t.username, t.password)
	}
	return t.base.RoundTrip(req)
}
