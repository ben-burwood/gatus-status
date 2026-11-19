package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadConfigFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	os.WriteFile(path, []byte("title: Acme Status\nlogo: https://acme.example/logo.png\ngatus:\n  url: https://gatus.example/\nendpoints:\n  - frontend\n  - core/api\n"), 0o644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Title != "Acme Status" || cfg.Logo != "https://acme.example/logo.png" || cfg.Gatus.URL != "https://gatus.example" {
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
	if !cfg.AnnouncementsEnabled() {
		t.Error("expected announcements to be enabled by default")
	}
}

func TestLoadConfigAnnouncements(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yml")
	os.WriteFile(path, []byte("gatus:\n  url: https://gatus.example\nendpoints:\n  - frontend\nannouncements: false\n"), 0o644)

	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.AnnouncementsEnabled() {
		t.Error("expected announcements to be disabled from the config file")
	}

	t.Setenv(AnnouncementsEnvVar, "true")
	if cfg, err = LoadConfig(path); err != nil {
		t.Fatal(err)
	}
	if !cfg.AnnouncementsEnabled() {
		t.Errorf("expected %s to override the config file", AnnouncementsEnvVar)
	}

	t.Setenv(AnnouncementsEnvVar, "nope")
	if _, err = LoadConfig(path); err == nil {
		t.Errorf("expected an invalid %s value to fail", AnnouncementsEnvVar)
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
