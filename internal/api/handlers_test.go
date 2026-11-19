package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gatus-status/internal/config"
	"gatus-status/internal/gatus"
)

const gatusStatuses = `[
	{"name":"frontend","group":"core","key":"core_frontend","results":[
		{"status":200,"hostname":"internal.example","success":false,"timestamp":"2026-11-19T10:00:00Z","errors":["secret"]},
		{"status":200,"hostname":"internal.example","success":true,"timestamp":"2026-11-19T10:01:00Z"}
	]},
	{"name":"api","group":"core","key":"core_api","results":[
		{"success":false,"timestamp":"2026-11-19T10:01:00Z"}
	]},
	{"name":"blog","key":"_blog","results":[]}
]`

func newGatus(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/endpoints/statuses":
			w.Write([]byte(gatusStatuses))
		case strings.HasSuffix(r.URL.Path, "/uptimes/24h"):
			w.Write([]byte("0.990000"))
		case strings.Contains(r.URL.Path, "/uptimes/"):
			w.Write([]byte("0.500000"))
		default:
			http.NotFound(w, r)
		}
	}))
}

func TestStatusesHandler(t *testing.T) {
	srv := newGatus(t)
	defer srv.Close()

	cfg := &config.Config{
		Gatus:     config.GatusConfig{URL: srv.URL},
		Endpoints: []string{"core/api", "frontend", "_blog", "missing"},
	}
	handler := StatusesHandler(cfg, gatus.NewClient(cfg.Gatus))

	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, "/api/statuses", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if body := rec.Body.String(); strings.Contains(body, "secret") {
		t.Fatalf("response leaked errors: %s", body)
	}

	var statuses []EndpointStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &statuses); err != nil {
		t.Fatal(err)
	}

	expected := []struct {
		name  string
		state State
	}{
		{"api", StateDown},
		{"frontend", StateUp},
		{"blog", StateUnknown},
		{"missing", StateUnknown},
	}
	if len(statuses) != len(expected) {
		t.Fatalf("expected %d statuses, got %d", len(expected), len(statuses))
	}
	for i, e := range expected {
		if statuses[i].Name != e.name || statuses[i].State != e.state {
			t.Errorf("status %d: expected %s/%s, got %s/%s", i, e.name, e.state, statuses[i].Name, statuses[i].State)
		}
	}

	if got := statuses[1].Uptime["24h"]; got != 0.99 {
		t.Errorf("expected 24h uptime 0.99, got %v", got)
	}
	if got := statuses[1].Uptime["30d"]; got != 0.5 {
		t.Errorf("expected 30d uptime 0.5, got %v", got)
	}
	if got := statuses[1].Hostname; got != "internal.example" {
		t.Errorf("expected hostname internal.example, got %q", got)
	}
	if len(statuses[1].History) != 2 {
		t.Errorf("expected 2 history points, got %d", len(statuses[1].History))
	}
}

func TestStatusesHandlerGatusDown(t *testing.T) {
	cfg := &config.Config{Gatus: config.GatusConfig{URL: "http://127.0.0.1:1"}, Endpoints: []string{"a"}}
	handler := StatusesHandler(cfg, gatus.NewClient(cfg.Gatus))

	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, "/api/statuses", nil))
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}
