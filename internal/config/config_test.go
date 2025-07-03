package config

import (
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/utils"
)

func TestRunInit_CreatesConfigFile(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	err := RunInit("./src", "", "", "", "", false)
	if err != nil {
		t.Fatalf("RunInit failed: %v", err)
	}

	exists, err := afero.Exists(AppFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Error checking if config file exists: %v", err)
	}
	if !exists {
		t.Error("Expected skelly.yaml to be created")
	}

	content, err := afero.ReadFile(AppFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "src_dir: ./src") {
		t.Errorf("Config should contain src_dir: ./src, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "*.tsx") {
		t.Errorf("Config should contain default file patterns, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "node_modules") {
		t.Errorf("Config should contain default exclude dirs, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "{component}-skeleton.{ext}") {
		t.Errorf("Config should contain default naming pattern, got: %s", contentStr)
	}
}

func TestRunInit_WithCustomParameters(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	err := RunInit("./components", "*.vue,*.svelte", "dist,build", "*.spec.*", "{name}-loading.{ext}", true)
	if err != nil {
		t.Fatalf("RunInit failed: %v", err)
	}

	content, err := afero.ReadFile(AppFs, "skelly.yaml")
	if err != nil {
		t.Fatalf("Failed to read created config file: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "src_dir: ./components") {
		t.Errorf("Config should contain custom src_dir, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "*.vue") {
		t.Errorf("Config should contain custom file patterns, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "dist") {
		t.Errorf("Config should contain custom exclude dirs, got: %s", contentStr)
	}
	if !strings.Contains(contentStr, "{name}-loading.{ext}") {
		t.Errorf("Config should contain custom naming pattern, got: %s", contentStr)
	}
}

func TestRunInit_ConfigAlreadyExists(t *testing.T) {
	originalFs := AppFs
	originalUtilsFs := utils.AppFs
	memFs := afero.NewMemMapFs()
	AppFs = memFs
	utils.AppFs = memFs
	defer func() { 
		AppFs = originalFs 
		utils.AppFs = originalUtilsFs
	}()

	afero.WriteFile(memFs, "skelly.yaml", []byte("existing config"), 0644)

	err := RunInit("./src", "", "", "", "", false)
	if err == nil {
		t.Error("RunInit should return error when config already exists")
		return
	}

	expectedMsg := "Config already exists, stopping."
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}