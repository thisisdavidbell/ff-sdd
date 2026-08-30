package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunScriptStructureAndErrorHandling(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve repo root: %v", err)
	}

	scriptPath := filepath.Join(repoRoot, "run.sh")
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("failed to read run.sh: %v", err)
	}

	content := string(data)

	if !strings.Contains(content, "set -euo pipefail") && !strings.Contains(content, "set -e") {
		t.Errorf("run.sh should contain strict error handling (set -euo pipefail or set -e)")
	}

	captureIdx := strings.Index(content, "go run ./cmd/capture-guysports")
	processIdx := strings.Index(content, "go run ./cmd/process")
	renderIdx := strings.Index(content, "go run ./cmd/render")

	if captureIdx == -1 {
		t.Errorf("run.sh missing go run ./cmd/capture-guysports invocation")
	}
	if processIdx == -1 {
		t.Errorf("run.sh missing go run ./cmd/process invocation")
	}
	if renderIdx == -1 {
		t.Errorf("run.sh missing go run ./cmd/render invocation")
	}

	if !(captureIdx < processIdx && processIdx < renderIdx) {
		t.Errorf("run.sh stages are not in sequential order: capture (%d) -> process (%d) -> render (%d)", captureIdx, processIdx, renderIdx)
	}
}
