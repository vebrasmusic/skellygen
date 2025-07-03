package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
	"github.com/vebrasmusic/skellygen/internal/generation"
	"github.com/vebrasmusic/skellygen/internal/naming"
	"github.com/vebrasmusic/skellygen/internal/utils"
	"github.com/vebrasmusic/skellygen/internal/validation"
)

func setupTestFileSystem() afero.Fs {
	memFs := afero.NewMemMapFs()
	
	config.AppFs = memFs
	utils.AppFs = memFs
	discovery.AppFs = memFs
	naming.AppFs = memFs
	validation.AppFs = memFs
	generation.AppFs = memFs
	
	return memFs
}

func restoreFileSystem() {
	osFs := afero.NewOsFs()
	config.AppFs = osFs
	utils.AppFs = osFs
	discovery.AppFs = osFs
	naming.AppFs = osFs
	validation.AppFs = osFs
	generation.AppFs = osFs
}

func TestInitCommand(t *testing.T) {
	memFs := setupTestFileSystem()
	defer restoreFileSystem()

	cmd := &cobra.Command{}
	cmd.AddCommand(initCmd)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"init", "--src-dir", "./test-src"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	exists, err := afero.Exists(memFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Error checking config file: %v", err)
	}
	if !exists {
		t.Error("Expected skelly.yaml to be created")
	}

	content, err := afero.ReadFile(memFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "src_dir: ./test-src") {
		t.Errorf("Config should contain custom src_dir, got: %s", contentStr)
	}
}

func TestInitCommandWithDefaults(t *testing.T) {
	memFs := setupTestFileSystem()
	defer restoreFileSystem()

	SrcDir = "./src"

	cmd := &cobra.Command{}
	cmd.AddCommand(initCmd)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	content, err := afero.ReadFile(memFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Error reading config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "src_dir: ./src") {
		t.Errorf("Config should contain default src_dir, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "*.tsx") {
		t.Errorf("Config should contain default file patterns, got: %s", contentStr)
	}
}

func TestInitCommandConfigExists(t *testing.T) {
	memFs := setupTestFileSystem()
	defer restoreFileSystem()

	afero.WriteFile(memFs, "skelly.yaml", []byte("existing"), 0644)

	cmd := &cobra.Command{}
	cmd.AddCommand(initCmd)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"init"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected init command to fail when config exists")
		return
	}

	expectedMsg := "Config already exists, stopping."
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message about existing config, got: %v", err)
	}
}

func TestGenCommand(t *testing.T) {
	memFs := setupTestFileSystem()
	defer restoreFileSystem()

	memFs.MkdirAll("/test/src/components", 0755)

	configContent := `
project:
  name: test-project
  version: 1.0.0
input:
  src_dir: /test/src
  file_patterns:
    - "*.tsx"
  exclude_dirs:
    - "node_modules"
  exclude_files:
    - "*.test.*"
output:
  naming_pattern: "{component}-skeleton.{ext}"
  preserve_structure: true
`
	afero.WriteFile(memFs, "skelly.yaml", []byte(configContent), 0644)

	sourceContent := `export const Button = () => <button>Click</button>;`
	afero.WriteFile(memFs, "/test/src/components/Button.tsx", []byte(sourceContent), 0644)

	cmd := &cobra.Command{}
	cmd.AddCommand(genCmd)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"gen"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("gen command failed: %v", err)
	}

	exists, err := afero.Exists(memFs, "/test/src/components/Button-skeleton.tsx")
	if err != nil {
		t.Fatalf("Error checking skeleton file: %v", err)
	}
	if !exists {
		t.Error("Expected skeleton file to be created")
	}

	skeletonContent, err := afero.ReadFile(memFs, "/test/src/components/Button-skeleton.tsx")
	if err != nil {
		t.Fatalf("Error reading skeleton file: %v", err)
	}

	if string(skeletonContent) != sourceContent {
		t.Errorf("Skeleton content doesn't match source")
	}
}

func TestGenCommandNoConfig(t *testing.T) {
	setupTestFileSystem()
	defer restoreFileSystem()

	cmd := &cobra.Command{}
	cmd.AddCommand(genCmd)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	cmd.SetArgs([]string{"gen"})

	err := cmd.Execute()
	if err == nil {
		t.Error("Expected gen command to fail when no config exists")
	}

	expectedMsg := "skelly.yaml not found. Please run 'skelly init' first to create a configuration file"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error message about missing config, got: %v", err)
	}
}