package unit

import (
	"testing"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
)

func TestDefaultConfigUses2026_27Season(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Season != "2026-27" {
		t.Fatalf("expected season 2026-27, got %q", cfg.Season)
	}
	if cfg.RawDir != "data/2026-27/raw" {
		t.Fatalf("expected raw dir data/2026-27/raw, got %q", cfg.RawDir)
	}
}
