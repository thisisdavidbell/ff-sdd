package capture

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
)

// WriteFixtureSnapshots writes fixture snapshots to a directory as append-only YAML files in per-manager directories.
func WriteFixtureSnapshots(outDir string, snapshots []models.ManagerSnapshot) error {
	for _, snapshot := range snapshots {
		teamName := snapshot.TeamName
		if teamName == "" {
			teamName = snapshot.ManagerName
		}
		if teamName == "" {
			teamName = fmt.Sprintf("team_%s", snapshot.ManagerID)
		}
		dir := filepath.Join(outDir, fmt.Sprintf("%s_%s", teamName, snapshot.ManagerID))
		filename := fmt.Sprintf("%s.yaml", strings.ReplaceAll(snapshot.CapturedAt, ":", "-"))
		if err := storage.WriteSnapshotFile(filepath.Join(dir, filename), snapshot); err != nil {
			return err
		}
	}
	return nil
}
