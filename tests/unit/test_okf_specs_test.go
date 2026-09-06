package unit

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

var requiredOKFFields = []string{
	"type",
	"title",
	"description",
	"tags",
	"status",
	"feature",
	"sdd_approach",
	"input_summary",
	"generated",
}

func TestSpecificationsConformToOKFProfile(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}

	specsRoot := filepath.Join(repoRoot, "specs")
	err = filepath.WalkDir(specsRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}

		name := filepath.Base(path)
		if name == "index.md" || name == "log.md" {
			return nil
		}

		metadata := readOKFFrontmatter(t, path)
		for _, field := range requiredOKFFields {
			if !valuePresent(metadata[field]) {
				t.Errorf("%s is missing required OKF field %q", path, field)
			}
		}
		if _, exists := metadata["feature_branch"]; exists {
			t.Errorf("%s contains redundant field %q", path, "feature_branch")
		}
		if status := stringValue(metadata["status"]); status != "draft" && status != "stable" && status != "deprecated" {
			t.Errorf("%s has invalid status %q", path, status)
		}
		if approach := stringValue(metadata["sdd_approach"]); approach != "streamlined" && approach != "full-speckit" {
			t.Errorf("%s has invalid sdd_approach %q", path, approach)
		}
		validateActorTimestamp(t, path, "generated", metadata["generated"])
		if verified, ok := metadata["verified"]; ok {
			validateVerification(t, path, verified)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("failed to walk specs: %v", err)
	}
}

func TestOKFIndexDeclaresVersionAndListsFeatures(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}
	indexPath := filepath.Join(repoRoot, "specs", "index.md")
	metadata := readOKFFrontmatter(t, indexPath)
	if version := stringValue(metadata["okf_version"]); version != "0.2" {
		t.Errorf("specs/index.md should declare okf_version 0.2, got %q", version)
	}
	if len(metadata) != 1 {
		t.Errorf("specs/index.md should only declare okf_version")
	}

	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("failed to read index: %v", err)
	}
	for _, feature := range []string{"001-", "002-", "003-", "004-", "005-", "006-", "007-", "008-", "009-", "010-", "011-", "012-", "013-"} {
		if !strings.Contains(string(content), "("+feature) {
			t.Errorf("specs/index.md is missing feature %s", feature)
		}
	}
}

func TestDevelopmentProcessDefinesSDDApprovalGate(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to resolve repository root: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(repoRoot, "DEVELOPMENT-PROCESS.md"))
	if err != nil {
		t.Fatalf("failed to read development process: %v", err)
	}
	process := string(content)
	for _, requirement := range []string{
		"## SDD Approval Gate",
		"explicit user approval",
		"status: draft",
		"status: stable",
		"Do not infer approval from silence",
	} {
		if !strings.Contains(process, requirement) {
			t.Errorf("DEVELOPMENT-PROCESS.md is missing approval-gate requirement %q", requirement)
		}
	}
}

func readOKFFrontmatter(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	parts := strings.SplitN(string(content), "---\n", 3)
	if len(parts) != 3 || parts[0] != "" {
		t.Fatalf("%s must begin with YAML frontmatter", path)
	}
	metadata := map[string]interface{}{}
	if err := yaml.Unmarshal([]byte(parts[1]), &metadata); err != nil {
		t.Fatalf("%s has invalid YAML frontmatter: %v", path, err)
	}
	return metadata
}

func validateVerification(t *testing.T, path string, value interface{}) {
	t.Helper()
	switch events := value.(type) {
	case map[string]interface{}:
		validateActorTimestamp(t, path, "verified", events)
	case []interface{}:
		for _, event := range events {
			mapping, ok := event.(map[string]interface{})
			if !ok {
				t.Errorf("%s has invalid verified event", path)
				continue
			}
			validateActorTimestamp(t, path, "verified", mapping)
		}
	default:
		t.Errorf("%s has invalid verified value", path)
	}
}

func validateActorTimestamp(t *testing.T, path, field string, value interface{}) {
	t.Helper()
	mapping, ok := value.(map[string]interface{})
	if !ok || strings.TrimSpace(stringValue(mapping["by"])) == "" || !valuePresent(mapping["at"]) {
		t.Errorf("%s has invalid %s actor and timestamp", path, field)
	}
}

func valuePresent(value interface{}) bool {
	if value == nil {
		return false
	}
	if text, ok := value.(string); ok {
		return strings.TrimSpace(text) != ""
	}
	if items, ok := value.([]interface{}); ok {
		return len(items) > 0
	}
	if mapping, ok := value.(map[string]interface{}); ok {
		return len(mapping) > 0
	}
	return strings.TrimSpace(fmt.Sprint(value)) != ""
}

func stringValue(value interface{}) string {
	return fmt.Sprint(value)
}
