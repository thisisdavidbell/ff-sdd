package unit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
)

func TestDefaultConfigUses2026_27Season(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Season != "2026-27" {
		t.Fatalf("expected season 2026-27, got %q", cfg.Season)
	}
	if cfg.Pages != 3 {
		t.Fatalf("expected 3 pages, got %d", cfg.Pages)
	}
	if cfg.RawDir != "data/2026-27/raw" {
		t.Fatalf("expected raw dir data/2026-27/raw, got %q", cfg.RawDir)
	}
}

func TestLoadConfigFromYAMLFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	data := []byte("season: \"2027-28\"\npages: 4\nsource_url: \"https://example.com/season.php\"\n")
	if err := os.WriteFile(configPath, data, 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("GUYSPORTS_CONFIG", configPath)
	cfg := config.Load()
	if cfg.Season != "2027-28" {
		t.Fatalf("expected season 2027-28, got %q", cfg.Season)
	}
	if cfg.Pages != 4 {
		t.Fatalf("expected 4 pages, got %d", cfg.Pages)
	}
	if cfg.RawDir != "data/2027-28/raw" {
		t.Fatalf("expected raw dir data/2027-28/raw, got %q", cfg.RawDir)
	}
}
