package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
)

func TestCaptureFlowCreatesHistoricalSnapshots(t *testing.T) {
	outDir := t.TempDir()
	base := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
	fixture := []models.ManagerSnapshot{
		{ManagerID: "m1", ManagerName: "Alpha_Mgr", TeamName: "Alpha_FC", CapturedAt: base.Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p2"}}},
		{ManagerID: "m2", ManagerName: "Beta_Mgr", TeamName: "Beta_FC", CapturedAt: base.Format(time.RFC3339), Source: "guysports", Players: []models.PlayerReference{{PlayerID: "p1"}, {PlayerID: "p3"}}},
	}

	if err := capture.WriteFixtureSnapshots(outDir, fixture); err != nil {
		t.Fatalf("WriteFixtureSnapshots returned error: %v", err)
	}

	loaded, err := storage.SnapshotFiles(outDir)
	if err != nil {
		t.Fatalf("SnapshotFiles returned error: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 loaded snapshots from nested directories, got %d", len(loaded))
	}

	// Verify per-manager directories were created
	m1Dir := filepath.Join(outDir, "Alpha_FC_m1")
	if stat, err := os.Stat(m1Dir); err != nil || !stat.IsDir() {
		t.Fatalf("expected directory %s to exist", m1Dir)
	}
}
