package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"gatus-status/internal/config"

	sdk "github.com/TwiN/gatus-sdk"
)

const (
	statusCacheTTL = 30 * time.Second
	refreshTimeout = 30 * time.Second
)

// uptimeDurations are the uptime windows fetched from Gatus for each endpoint
var uptimeDurations = []string{"24h", "7d", "30d"}

// State is the simplified public state of an endpoint
type State string

const (
	StateUp      State = "up"
	StateDown    State = "down"
	StateUnknown State = "unknown"
)

// HistoryPoint is a single public check result - deliberately excludes errors and conditions
type HistoryPoint struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// EndpointStatus is the public status of a single configured endpoint
type EndpointStatus struct {
	Name        string             `json:"name"`
	Group       string             `json:"group,omitempty"`
	State       State              `json:"state"`
	Hostname    string             `json:"hostname,omitempty"`
	LastChecked *time.Time         `json:"lastChecked,omitempty"`
	Uptime      map[string]float64 `json:"uptime"`
	History     []HistoryPoint     `json:"history"`
}

// ConfigHandler returns the public page configuration
func ConfigHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"title": cfg.Title})
	}
}

// StatusesHandler returns the public status of every configured endpoint, in configured order
func StatusesHandler(cfg *config.Config, client *sdk.Client) http.HandlerFunc {
	var mu sync.Mutex
	var cached []EndpointStatus
	var expiresAt time.Time

	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if cached == nil || time.Now().After(expiresAt) {
			// Detached from the request context, as the result is cached and shared with other visitors
			ctx, cancel := context.WithTimeout(context.Background(), refreshTimeout)
			statuses, err := fetchStatuses(ctx, cfg, client)
			cancel()
			if err != nil {
				log.Printf("[api.StatusesHandler] Failed to fetch statuses from Gatus: %v", err)
				http.Error(w, "failed to fetch statuses", http.StatusBadGateway)
				return
			}
			cached, expiresAt = statuses, time.Now().Add(statusCacheTTL)
		}

		writeJSON(w, cached)
	}
}

func fetchStatuses(ctx context.Context, cfg *config.Config, client *sdk.Client) ([]EndpointStatus, error) {
	all, err := client.GetAllEndpointStatuses(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]EndpointStatus, len(cfg.Endpoints))
	var mu sync.Mutex
	var wg sync.WaitGroup
	for i, name := range cfg.Endpoints {
		gs := findEndpoint(all, name)
		if gs == nil {
			log.Printf("[api.fetchStatuses] Endpoint %q not found in Gatus", name)
			results[i] = EndpointStatus{Name: name, State: StateUnknown, Uptime: map[string]float64{}, History: []HistoryPoint{}}
			continue
		}
		results[i] = toEndpointStatus(gs)

		for _, duration := range uptimeDurations {
			wg.Add(1)
			go func() {
				defer wg.Done()
				uptime, err := client.GetEndpointUptime(ctx, gs.Key, duration)
				if err != nil {
					log.Printf("[api.fetchStatuses] Failed to fetch %s uptime for %s: %v", duration, gs.Key, err)
					return
				}
				mu.Lock()
				results[i].Uptime[duration] = uptime
				mu.Unlock()
			}()
		}
	}
	wg.Wait()
	return results, nil
}

// findEndpoint matches a configured name against a Gatus endpoint's name, "group/name" or key
func findEndpoint(all []sdk.EndpointStatus, name string) *sdk.EndpointStatus {
	for i := range all {
		s := &all[i]
		if s.Name == name || s.Key == name || (s.Group != "" && s.Group+"/"+s.Name == name) {
			return s
		}
	}
	return nil
}

func toEndpointStatus(gs *sdk.EndpointStatus) EndpointStatus {
	status := EndpointStatus{
		Name:    gs.Name,
		Group:   gs.Group,
		State:   StateUnknown,
		Uptime:  map[string]float64{},
		History: make([]HistoryPoint, 0, len(gs.Results)),
	}
	for _, result := range gs.Results {
		status.History = append(status.History, HistoryPoint{Success: result.Success, Timestamp: result.Timestamp})
	}
	if n := len(gs.Results); n > 0 {
		last := gs.Results[n-1]
		status.Hostname = last.Hostname
		status.LastChecked = &last.Timestamp
		if last.Success {
			status.State = StateUp
		} else {
			status.State = StateDown
		}
	}
	return status
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
