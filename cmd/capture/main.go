package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/thisisdavidbell/ff-sdd/internal/capture"
	"github.com/thisisdavidbell/ff-sdd/internal/config"
	"github.com/thisisdavidbell/ff-sdd/internal/models"
)

func main() {
	cfg := config.Load()
	source := flag.String("source", "guysports", "source name")
	output := flag.String("output", cfg.RawDir, "output directory")
	flag.Parse()

	if _, err := os.Stat(*output); err != nil {
		if err := os.MkdirAll(*output, 0o755); err != nil {
			panic(err)
		}
	}

	if *source != "guysports" {
		fmt.Fprintf(os.Stderr, "unsupported source: %s\n", *source)
		os.Exit(2)
	}

	resp, err := http.Get(cfg.SourceURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "unexpected HTTP status: %s\n", resp.Status)
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read body failed: %v\n", err)
		os.Exit(1)
	}
	html := string(body)
	if !strings.Contains(html, "GuySports") && !strings.Contains(html, "Season Table") {
		fmt.Fprintf(os.Stderr, "unexpected live page payload from Guy Sports\n")
		os.Exit(1)
	}

	matches := capture.ExtractManagersFromHTML(html)
	if len(matches) == 0 {
		fmt.Fprintf(os.Stderr, "no managers found in live response\n")
		os.Exit(1)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	count := 0
	for _, manager := range matches {
		players, err := capture.FetchAndParseManagerDetail(manager.DetailURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch manager detail failed for %s: %v\n", manager.Name, err)
			continue
		}
		snapshot := models.ManagerSnapshot{
			ManagerID:   manager.ID,
			ManagerName: manager.Name,
			CapturedAt:  now,
			Source:      *source,
			Players:     players,
			RawURL:      manager.DetailURL,
		}
		path := filepath.Join(*output, fmt.Sprintf("%s-%s.yaml", sanitize(manager.ID), strings.ReplaceAll(now, ":", "-")))
		if err := capture.WriteSnapshot(path, snapshot); err != nil {
			fmt.Fprintf(os.Stderr, "write snapshot failed: %v\n", err)
			os.Exit(1)
		}
		count++
	}
	fmt.Printf("captured %d managers into %s\n", count, *output)
}

func sanitize(s string) string {
	return strings.NewReplacer(" ", "-", "/", "-", "?", "-", "&", "-", "=", "-").Replace(s)
}
