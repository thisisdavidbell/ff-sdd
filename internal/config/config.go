package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of config.yaml.
type FileConfig struct {
	Season    string `yaml:"season"`
	Pages     int    `yaml:"pages"`
	SourceURL string `yaml:"source_url"`
}

// Config holds runtime configuration and derived paths.
type Config struct {
	Season     string
	Pages      int
	RawDir     string
	ProcessDir string
	RenderFile string
	SourceURL  string
}

// DefaultConfig returns the fallback configuration.
func DefaultConfig() Config {
	return Config{
		Season:     "2026-27",
		Pages:      3,
		RawDir:     filepath.Join("data", "2026-27", "raw"),
		ProcessDir: filepath.Join("data", "2026-27", "processed"),
		RenderFile: filepath.Join("docs", "index.html"),
		SourceURL:  "https://www.guysports.co.uk/guysports/season.php",
	}
}

// Load reads configuration from config.yaml or GUYSPORTS_CONFIG, falling back to defaults.
func Load() Config {
	cfg := DefaultConfig()

	configPath := "config.yaml"
	if envPath := os.Getenv("GUYSPORTS_CONFIG"); envPath != "" {
		configPath = envPath
	}

	if data, err := os.ReadFile(configPath); err == nil {
		var fc FileConfig
		if err := yaml.Unmarshal(data, &fc); err == nil {
			if fc.Season != "" {
				cfg.Season = fc.Season
			}
			if fc.Pages > 0 {
				cfg.Pages = fc.Pages
			}
			if fc.SourceURL != "" {
				cfg.SourceURL = fc.SourceURL
			}
		}
	}

	if v := os.Getenv("GUYSPORTS_SEASON"); v != "" {
		cfg.Season = v
	}

	cfg.RawDir = filepath.Join("data", cfg.Season, "raw")
	cfg.ProcessDir = filepath.Join("data", cfg.Season, "processed")
	cfg.RenderFile = filepath.Join("docs", "index.html")

	return cfg
}
