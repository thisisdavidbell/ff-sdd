package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/processing"
	"github.com/thisisdavidbell/ff-sdd/internal/storage"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := config.Load()
	inputDir := cfg.RawDir
	currentOutDir := filepath.Join(cfg.ProcessDir, "player-ownership")
	changesOutDir := filepath.Join(cfg.ProcessDir, "manager-changes")
	if err := os.MkdirAll(currentOutDir, 0o755); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(changesOutDir, 0o755); err != nil {
		panic(err)
	}

	snapshots, err := storage.SnapshotFiles(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "process failed: %v\n", err)
		os.Exit(1)
	}
	current, history, err := processing.BuildOwnership(snapshots)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ownership build failed: %v\n", err)
		os.Exit(1)
	}
	for id, rec := range current {
		data, _ := yaml.Marshal(models.PlayerOwnershipRecord{PlayerID: rec.PlayerID, PlayerName: rec.PlayerName, ManagerCount: rec.ManagerCount})
		if err := os.WriteFile(filepath.Join(currentOutDir, id+".yaml"), data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write ownership failed: %v\n", err)
			os.Exit(1)
		}
	}
	changes, err := processing.BuildChangeSummaries(snapshots)
	if err != nil {
		fmt.Fprintf(os.Stderr, "change summary failed: %v\n", err)
		os.Exit(1)
	}
	for _, summary := range changes {
		data, _ := yaml.Marshal(summary)
		if err := os.WriteFile(filepath.Join(changesOutDir, summary.ManagerID+".yaml"), data, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "write changes failed: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("processed %d snapshot records\n", len(snapshots))
	_ = history
}
