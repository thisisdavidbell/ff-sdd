package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

func TestCaptureFlowCreatesHistoricalSnapshots(t *testing.T) {
	outDir := t.TempDir()
	base := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	fixture := []models.ManagerSnapshot{
		{ManagerID: "m1", ManagerName: "A", CapturedAt: base.Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p2"}}},
		{ManagerID: "m2", ManagerName: "B", CapturedAt: base.Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
	}

	if err := capture.WriteFixtureSnapshots(outDir, fixture); err != nil {
		t.Fatalf("WriteFixtureSnapshots returned error: %v", err)
	}

	files, err := filepath.Glob(filepath.Join(outDir, "*.yaml"))
	if err != nil {
		t.Fatalf("Glob returned error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 snapshot files, got %d", len(files))
	}

	for _, p := range files {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("snapshot file missing: %v", err)
		}
	}
}
