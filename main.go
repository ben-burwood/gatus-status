package main

import (
	"gatus-status/internal/api"
	"gatus-status/internal/config"
	"gatus-status/internal/gatus"
	"log"
	"net/http"
	"os"
)

const (
	ConfigPathEnvVar = "GATUS_STATUS_CONFIG_PATH"
)

func main() {
	cfg, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	client := gatus.NewClient(cfg.Gatus)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/config", api.ConfigHandler(cfg))
	mux.HandleFunc("GET /api/statuses", api.StatusesHandler(cfg, client))
	mux.HandleFunc("GET /api/announcements", api.AnnouncementsHandler(cfg, client))

	// Serve Static Frontend
	mux.Handle("/", http.FileServer(http.Dir("./frontend/dist")))

	log.Printf("Listening on :8080 (gatus=%s, endpoints=%d)", cfg.Gatus.URL, len(cfg.Endpoints))
	log.Fatal(http.ListenAndServe("[::]:8080", api.CORSMiddleware(mux)))
}

// loadConfiguration loads the configuration from the path specified in the GATUS_STATUS_CONFIG_PATH environment variable
func loadConfiguration() (*config.Config, error) {
	configPath := os.Getenv(ConfigPathEnvVar)
	return config.LoadConfig(configPath)
}
