package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
	"github.com/thisisdavidbell/ff-sdd/internal/render"
	"gopkg.in/yaml.v3"
)

func main() {
	cfg := config.Load()
	inputDir := filepath.Join(cfg.ProcessDir, "player-ownership")
	changesDir := filepath.Join(cfg.ProcessDir, "manager-changes")
	outputFile := cfg.RenderFile

	current := map[string]models.PlayerOwnershipRecord{}
	files, err := os.ReadDir(inputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to read processed data: %v\n", err)
		os.Exit(1)
	}
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
			continue
		}
		path := filepath.Join(inputDir, file.Name())
		var rec models.PlayerOwnershipRecord
		data, _ := os.ReadFile(path)
		_ = yaml.Unmarshal(data, &rec)
		if rec.PlayerID != "" {
			current[rec.PlayerID] = rec
		}
	}

	history := map[string][]models.PlayerOwnershipHistoryEntry{}
	changes := []models.ManagerChangeSummary{}
	changeFiles, err := os.ReadDir(changesDir)
	if err == nil {
		for _, file := range changeFiles {
			if file.IsDir() || filepath.Ext(file.Name()) != ".yaml" {
				continue
			}
			path := filepath.Join(changesDir, file.Name())
			var summary models.ManagerChangeSummary
			data, _ := os.ReadFile(path)
			_ = yaml.Unmarshal(data, &summary)
			if summary.ManagerID != "" {
				changes = append(changes, summary)
			}
		}
	}
	if html, err := render.BuildHTML(current, history, changes); err != nil {
		fmt.Fprintf(os.Stderr, "render failed: %v\n", err)
		os.Exit(1)
	} else if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "unable to prepare output dir: %v\n", err)
		os.Exit(1)
	} else if err := os.WriteFile(outputFile, []byte(html), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "unable to write render output: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("rendered report to %s\n", outputFile)
}
