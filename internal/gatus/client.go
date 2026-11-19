package gatus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gatus-status/internal/config"
)

// resultsPerEndpoint is how many recent results Gatus returns per endpoint.
// Note Gatus' page/pageSize paginate each endpoint's results - all endpoints are always returned.
const resultsPerEndpoint = 50

// EndpointStatus mirrors the subset of Gatus' endpoint.Status that we use
type EndpointStatus struct {
	Name    string   `json:"name"`
	Group   string   `json:"group"`
	Key     string   `json:"key"`
	Results []Result `json:"results"`
}

// Result mirrors the subset of Gatus' endpoint.Result that we use
type Result struct {
	Success   bool          `json:"success"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
}

type Client struct {
	cfg  config.GatusConfig
	http *http.Client
}

func NewClient(cfg config.GatusConfig) *Client {
	return &Client{cfg: cfg, http: &http.Client{Timeout: 10 * time.Second}}
}

// AllStatuses fetches every endpoint status from Gatus, with its most recent results
func (c *Client) AllStatuses() ([]EndpointStatus, error) {
	var statuses []EndpointStatus
	path := fmt.Sprintf("/api/v1/endpoints/statuses?pageSize=%d", resultsPerEndpoint)
	if err := c.getJSON(path, &statuses); err != nil {
		return nil, err
	}
	return statuses, nil
}

// Uptime fetches the uptime (0-1) for an endpoint key over a duration (1h, 24h, 7d, 30d)
func (c *Client) Uptime(key, duration string) (float64, error) {
	body, err := c.get(fmt.Sprintf("/api/v1/endpoints/%s/uptimes/%s", url.PathEscape(key), duration))
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(strings.TrimSpace(string(body)), 64)
}

func (c *Client) getJSON(path string, v any) error {
	body, err := c.get(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

func (c *Client) get(path string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.cfg.URL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.cfg.Username != "" || c.cfg.Password != "" {
		req.SetBasicAuth(c.cfg.Username, c.cfg.Password)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gatus %s returned %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return body, nil
}
