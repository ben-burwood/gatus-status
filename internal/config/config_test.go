package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	os.WriteFile(path, []byte("title: Acme Status\ngatus:\n  url: https://gatus.example/\nendpoints:\n  - frontend\n  - core/api\n"), 0o644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Title != "Acme Status" || cfg.Gatus.URL != "https://gatus.example" {
		t.Errorf("unexpected config: %+v", cfg)
	}
	if !reflect.DeepEqual(cfg.Endpoints, []string{"frontend", "core/api"}) {
		t.Errorf("unexpected endpoints: %v", cfg.Endpoints)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	t.Chdir(t.TempDir()) // ensure the default config file isn't picked up
	t.Setenv(GatusURLEnvVar, "http://gatus:8080")
	t.Setenv(EndpointsEnvVar, " frontend, core/api ,,")

	cfg, err := LoadConfig("")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(cfg.Endpoints, []string{"frontend", "core/api"}) {
		t.Errorf("unexpected endpoints: %v", cfg.Endpoints)
	}
	if cfg.Title != DefaultTitle {
		t.Errorf("expected default title, got %q", cfg.Title)
	}
}

func TestLoadConfigValidation(t *testing.T) {
	t.Chdir(t.TempDir())

	if _, err := LoadConfig(""); err != ErrNoGatusURL {
		t.Errorf("expected ErrNoGatusURL, got %v", err)
	}
	t.Setenv(GatusURLEnvVar, "http://gatus:8080")
	if _, err := LoadConfig(""); err != ErrNoEndpointsInConfig {
		t.Errorf("expected ErrNoEndpointsInConfig, got %v", err)
	}
}
