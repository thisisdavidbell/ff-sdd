package capture

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
)

// WriteFixtureSnapshots writes fixture snapshots to a directory as append-only YAML files.
func WriteFixtureSnapshots(outDir string, snapshots []models.ManagerSnapshot) error {
	for _, snapshot := range snapshots {
		filename := fmt.Sprintf("%s-%s.yaml", snapshot.ManagerID, strings.ReplaceAll(snapshot.CapturedAt, ":", "-"))
		if err := storage.WriteSnapshotFile(filepath.Join(outDir, filename), snapshot); err != nil {
			return err
		}
	}
	return nil
}
