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

// SnapshotFiles reads all YAML files in a directory and returns their parsed snapshots.
func SnapshotFiles(dir string) ([]models.ManagerSnapshot, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []models.ManagerSnapshot
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		snap, err := ReadSnapshotFile(path)
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", path, err)
		}
		out = append(out, snap)
	}
	return out, nil
}
