package gatus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gatus-status/internal/config"
)

func TestNewClientBasicAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if username, password, ok := r.BasicAuth(); !ok || username != "user" || password != "pass" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	client := NewClient(config.GatusConfig{URL: srv.URL, Username: "user", Password: "pass"})
	if _, err := client.GetAllEndpointStatuses(context.Background()); err != nil {
		t.Fatalf("expected basic auth to succeed, got %v", err)
	}

	client = NewClient(config.GatusConfig{URL: srv.URL})
	if _, err := client.GetAllEndpointStatuses(context.Background()); err == nil {
		t.Fatal("expected request without basic auth to fail")
	}
}

func TestNewClientBasicAuthNotSentToOtherHosts(t *testing.T) {
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, _, ok := r.BasicAuth(); ok {
			t.Error("credentials were sent to a redirected host")
		}
		w.Write([]byte(`[]`))
	}))
	defer other.Close()
	gatus := httptest.NewServer(http.RedirectHandler(other.URL+"/api/v1/endpoints/statuses", http.StatusFound))
	defer gatus.Close()

	client := NewClient(config.GatusConfig{URL: gatus.URL, Username: "user", Password: "pass"})
	if _, err := client.GetAllEndpointStatuses(context.Background()); err != nil {
		t.Fatal(err)
	}
}
