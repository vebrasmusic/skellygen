package generation

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
	"github.com/vebrasmusic/skellygen/internal/naming"
	"github.com/vebrasmusic/skellygen/internal/utils"
	"github.com/vebrasmusic/skellygen/internal/validation"
)

func TestGetDirectories_ConfigExists(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	configContent := `
project:
  name: test-project
  version: 1.0.0
input:
  src_dir: ./src
  file_patterns:
    - "*.tsx"
    - "*.ts"
  exclude_dirs:
    - "node_modules"
  exclude_files:
    - "*.test.*"
output:
  naming_pattern: "{component}-skeleton.{ext}"
  preserve_structure: true
`
	afero.WriteFile(AppFs, "skelly.yaml", []byte(configContent), 0644)

	cfg, err := getDirectories()
	if err != nil {
		t.Fatalf("getDirectories() error = %v", err)
	}

	if cfg.Project.Name != "test-project" {
		t.Errorf("Expected project name 'test-project', got '%s'", cfg.Project.Name)
	}
	if cfg.Input.SrcDir != "./src" {
		t.Errorf("Expected src_dir './src', got '%s'", cfg.Input.SrcDir)
	}
	if len(cfg.Input.FilePatterns) != 2 {
		t.Errorf("Expected 2 file patterns, got %d", len(cfg.Input.FilePatterns))
	}
}

func TestGetDirectories_ConfigNotExists(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	_, err := getDirectories()
	if err == nil {
		t.Error("getDirectories() expected error when config doesn't exist")
	}

	expectedMsg := "skelly.yaml not found. Please run 'skelly init' first to create a configuration file"
	if err.Error() != expectedMsg {
		t.Errorf("getDirectories() error = %v, want %v", err.Error(), expectedMsg)
	}
}

func TestParseInputFile_Integration(t *testing.T) {
	originalAppFs := AppFs
	originalConfigFs := config.AppFs
	originalUtilsFs := utils.AppFs
	originalDiscoveryFs := discovery.AppFs
	originalNamingFs := naming.AppFs
	originalValidationFs := validation.AppFs

	memFs := afero.NewMemMapFs()
	AppFs = memFs
	config.AppFs = memFs
	utils.AppFs = memFs
	discovery.AppFs = memFs
	naming.AppFs = memFs
	validation.AppFs = memFs

	defer func() {
		AppFs = originalAppFs
		config.AppFs = originalConfigFs
		utils.AppFs = originalUtilsFs
		discovery.AppFs = originalDiscoveryFs
		naming.AppFs = originalNamingFs
		validation.AppFs = originalValidationFs
	}()

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

	sourceContent := `import React from 'react';

export const Button = () => {
  return <button>Click me</button>;
};`
	afero.WriteFile(memFs, "/test/src/components/Button.tsx", []byte(sourceContent), 0644)

	err := ParseInputFile()
	if err != nil {
		t.Fatalf("ParseInputFile() error = %v", err)
	}

	exists, err := afero.Exists(memFs, "/test/src/components/Button-skeleton.tsx")
	if err != nil {
		t.Fatalf("Error checking if skeleton file exists: %v", err)
	}
	if !exists {
		t.Error("Expected skeleton file to be created")
	}

	skeletonContent, err := afero.ReadFile(memFs, "/test/src/components/Button-skeleton.tsx")
	if err != nil {
		t.Fatalf("Error reading skeleton file: %v", err)
	}

	if string(skeletonContent) != sourceContent {
		t.Errorf("Skeleton content doesn't match source content")
	}
}

func TestGenerateSkeleton(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	sourceContent := "const test = 'hello';"
	AppFs.MkdirAll("/test", 0755)
	afero.WriteFile(AppFs, "/test/source.js", []byte(sourceContent), 0644)

	file := discovery.FileInfo{
		Path:         "/test/source.js",
		RelativePath: "source.js",
		Name:         "source",
		Extension:    "js",
	}

	err := generateSkeleton(file, "/test/source-skeleton.js")
	if err != nil {
		t.Fatalf("generateSkeleton() error = %v", err)
	}

	skeletonContent, err := afero.ReadFile(AppFs, "/test/source-skeleton.js")
	if err != nil {
		t.Fatalf("Error reading skeleton file: %v", err)
	}

	if string(skeletonContent) != sourceContent {
		t.Errorf("generateSkeleton() content = %v, want %v", string(skeletonContent), sourceContent)
	}
}