package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigurationFilePath is the default path that will be used to search for the configuration file
	// if a custom path isn't configured through the GATUS_STATUS_CONFIG_PATH environment variable
	DefaultConfigurationFilePath = "config/config.yml"

	// EndpointsEnvVar is the environment variable that can be used to configure endpoints as a comma-separated list
	EndpointsEnvVar = "ENDPOINTS"

	// GatusURLEnvVar is the environment variable for the base URL of the Gatus instance
	GatusURLEnvVar = "GATUS_URL"
	// GatusUsernameEnvVar and GatusPasswordEnvVar configure optional Basic Auth for the Gatus API
	GatusUsernameEnvVar = "GATUS_USERNAME"
	GatusPasswordEnvVar = "GATUS_PASSWORD"

	// TitleEnvVar is the environment variable for the page title
	TitleEnvVar = "TITLE"
	// LogoEnvVar is the environment variable for the URL of the logo shown above the title
	LogoEnvVar = "LOGO"

	// AnnouncementsEnvVar is the environment variable to show or hide the announcements configured in Gatus
	AnnouncementsEnvVar = "ANNOUNCEMENTS"

	DefaultTitle = "Status"
)

var (
	// ErrNoEndpointsInConfig is an error returned when no endpoints are configured
	ErrNoEndpointsInConfig = errors.New("no endpoints configured")

	// ErrNoGatusURL is an error returned when no Gatus URL is configured
	ErrNoGatusURL = errors.New("no gatus url configured")
)

type GatusConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Config struct {
	Title     string      `yaml:"title"`
	Logo      string      `yaml:"logo"` // URL of the logo shown above the title
	Gatus     GatusConfig `yaml:"gatus"`
	Endpoints []string    `yaml:"endpoints"`

	// Announcements shows the announcements configured in Gatus. Enabled by default if nil.
	Announcements *bool `yaml:"announcements"`
}

// AnnouncementsEnabled returns whether the announcements configured in Gatus should be shown
func (c *Config) AnnouncementsEnabled() bool {
	return c.Announcements == nil || *c.Announcements
}

// LoadConfig loads the configuration from the specified path (or the default path) if it exists,
// then applies any environment variable overrides. ENDPOINTS takes precedence over the config file endpoints.
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	for _, configurationPath := range []string{configPath, DefaultConfigurationFilePath} {
		if len(configurationPath) == 0 {
			continue
		}
		data, err := os.ReadFile(configurationPath)
		if err != nil {
			continue
		}
		log.Printf("[config.LoadConfig] Reading configuration from configFile=%s", configurationPath)
		if err = yaml.Unmarshal([]byte(os.ExpandEnv(string(data))), config); err != nil {
			return nil, fmt.Errorf("error parsing config %s: %w", configurationPath, err)
		}
		break
	}

	if envEndpoints, ok := os.LookupEnv(EndpointsEnvVar); ok {
		log.Printf("[config.LoadConfig] Reading endpoints from %s environment variable", EndpointsEnvVar)
		config.Endpoints = parseEndpointsEnv(envEndpoints)
	}
	if v := os.Getenv(GatusURLEnvVar); v != "" {
		config.Gatus.URL = v
	}
	if v := os.Getenv(GatusUsernameEnvVar); v != "" {
		config.Gatus.Username = v
	}
	if v := os.Getenv(GatusPasswordEnvVar); v != "" {
		config.Gatus.Password = v
	}
	if v := os.Getenv(TitleEnvVar); v != "" {
		config.Title = v
	}
	if v := os.Getenv(LogoEnvVar); v != "" {
		config.Logo = v
	}
	if v := os.Getenv(AnnouncementsEnvVar); v != "" {
		show, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("invalid %s value %q: %w", AnnouncementsEnvVar, v, err)
		}
		config.Announcements = &show
	}

	return config, config.validate()
}

func (c *Config) validate() error {
	c.Gatus.URL = strings.TrimRight(strings.TrimSpace(c.Gatus.URL), "/")
	if c.Gatus.URL == "" {
		return ErrNoGatusURL
	}
	if len(c.Endpoints) == 0 {
		return ErrNoEndpointsInConfig
	}
	if c.Title == "" {
		c.Title = DefaultTitle
	}
	return nil
}

// parseEndpointsEnv parses a comma-separated list of endpoint names
func parseEndpointsEnv(value string) []string {
	var endpoints []string
	for _, endpoint := range strings.Split(value, ",") {
		if endpoint = strings.TrimSpace(endpoint); len(endpoint) > 0 {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}
