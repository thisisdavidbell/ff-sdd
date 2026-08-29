package config

import (
	"os"
	"path/filepath"
)

// Config holds the project defaults for season-aware data locations.
type Config struct {
	Season    string
	RawDir    string
	ProcessDir string
	RenderFile string
	SourceURL string
}

// DefaultConfig returns the season-aware layout defaults.
func DefaultConfig() Config {
	return Config{
		Season:     "2026-27",
		RawDir:     filepath.Join("data", "2026-27", "raw"),
		ProcessDir: filepath.Join("data", "2026-27", "processed"),
		RenderFile: filepath.Join("docs", "index.html"),
		SourceURL:  "https://www.guysports.co.uk/guysports/season.php?page=1",
	}
}

// Load reads environment overrides for the season-aware config.
func Load() Config {
	cfg := DefaultConfig()
	if v := os.Getenv("GUYSPORTS_SEASON"); v != "" {
		cfg.Season = v
		cfg.RawDir = filepath.Join("data", v, "raw")
		cfg.ProcessDir = filepath.Join("data", v, "processed")
	}
	return cfg
}
