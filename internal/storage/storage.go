package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

// WriteSnapshotFile writes one manager snapshot as YAML.
func WriteSnapshotFile(path string, snapshot models.ManagerSnapshot) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(snapshot)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ReadSnapshotFile reads one YAML snapshot file.
func ReadSnapshotFile(path string) (models.ManagerSnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return models.ManagerSnapshot{}, err
	}
	var snapshot models.ManagerSnapshot
	if err := yaml.Unmarshal(data, &snapshot); err != nil {
		return models.ManagerSnapshot{}, err
	}
	return snapshot, nil
}

// SnapshotFiles reads all YAML files in a directory and any subdirectories, returning their parsed snapshots.
func SnapshotFiles(dir string) ([]models.ManagerSnapshot, error) {
	var out []models.ManagerSnapshot
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext != ".yaml" && ext != ".yml" {
			return nil
		}
		snap, readErr := ReadSnapshotFile(path)
		if readErr != nil {
			return fmt.Errorf("read %s: %w", path, readErr)
		}
		if snap.ManagerID != "" {
			out = append(out, snap)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
