package gatus

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"gatus-status/internal/config"

	sdk "github.com/TwiN/gatus-sdk"
)

const timeout = 10 * time.Second

// Client wraps the Gatus SDK client, adding the API calls the SDK does not cover
type Client struct {
	*sdk.Client
	baseURL    string
	httpClient *http.Client
}

// Announcement is a system-wide announcement configured in Gatus
type Announcement struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Archived  bool      `json:"archived,omitempty"`
}

// NewClient creates a Gatus client, adding Basic Auth if configured
func NewClient(cfg config.GatusConfig) *Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConnsPerHost = 10

	var roundTripper http.RoundTripper = transport
	if cfg.Username != "" || cfg.Password != "" {
		baseURL, _ := url.Parse(cfg.URL)
		roundTripper = &basicAuthTransport{username: cfg.Username, password: cfg.Password, host: baseURL.Host, base: transport}
	}

	httpClient := &http.Client{Timeout: timeout, Transport: roundTripper}
	return &Client{
		Client:     sdk.NewClient(cfg.URL, sdk.WithHTTPClient(httpClient)),
		baseURL:    cfg.URL,
		httpClient: httpClient,
	}
}

// GetAnnouncements returns the announcements from Gatus' /api/v1/config.
// Gatus versions without announcement support return none.
func (c *Client) GetAnnouncements(ctx context.Context) ([]Announcement, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/config", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d from %s", resp.StatusCode, req.URL.Path)
	}

	var body struct {
		Announcements []Announcement `json:"announcements"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Announcements == nil {
		body.Announcements = []Announcement{}
	}
	return body.Announcements, nil
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
