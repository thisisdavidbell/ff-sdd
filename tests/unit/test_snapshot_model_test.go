package unit

import (
	"testing"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
	"github.com/thisisdavidbell/ff-sdd/internal/validation"
)

func TestManagerSnapshotValidationAndYAMLRoundTrip(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	snapshot := models.ManagerSnapshot{
		ManagerID:   "mgr-1",
		ManagerName: "Alpha",
		CapturedAt:  now,
		Source:      "guysports",
		Players: []models.PlayerReference{
			{PlayerID: "p1", Name: "Haaland", Position: "FWD", TeamName: "Man City"},
			{PlayerID: "p2", Name: "Salah", Position: "MID", TeamName: "Liverpool"},
		},
	}

	if err := validation.ValidateSnapshot(snapshot); err != nil {
		t.Fatalf("ValidateSnapshot returned error: %v", err)
	}

	path := t.TempDir() + "/snapshot.yaml"
	if err := storage.WriteSnapshotFile(path, snapshot); err != nil {
		t.Fatalf("WriteSnapshotFile returned error: %v", err)
	}

	loaded, err := storage.ReadSnapshotFile(path)
	if err != nil {
		t.Fatalf("ReadSnapshotFile returned error: %v", err)
	}

	if loaded.ManagerID != snapshot.ManagerID || len(loaded.Players) != 2 {
		t.Fatalf("round trip mismatch: %#v", loaded)
	}
}
