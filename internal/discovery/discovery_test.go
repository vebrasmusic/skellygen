package discovery

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/config"
)

func TestFindFiles(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	AppFs.MkdirAll("/test/src/components", 0755)
	AppFs.MkdirAll("/test/src/utils", 0755)
	AppFs.MkdirAll("/test/src/node_modules", 0755)

	afero.WriteFile(AppFs, "/test/src/components/Button.tsx", []byte("button content"), 0644)
	afero.WriteFile(AppFs, "/test/src/components/Card.jsx", []byte("card content"), 0644)
	afero.WriteFile(AppFs, "/test/src/utils/helpers.ts", []byte("helpers content"), 0644)
	afero.WriteFile(AppFs, "/test/src/utils/helpers.test.ts", []byte("test content"), 0644)
	afero.WriteFile(AppFs, "/test/src/node_modules/package.js", []byte("package content"), 0644)
	afero.WriteFile(AppFs, "/test/src/README.md", []byte("readme content"), 0644)

	cfg := &config.Config{
		Input: config.Input{
			SrcDir:       "/test/src",
			FilePatterns: []string{"*.tsx", "*.ts", "*.jsx"},
			ExcludeDirs:  []string{"node_modules"},
			ExcludeFiles: []string{"*.test.*"},
		},
	}

	files, err := FindFiles(cfg)
	if err != nil {
		t.Fatalf("FindFiles() error = %v", err)
	}

	expectedFiles := map[string]bool{
		"/test/src/components/Button.tsx": true,
		"/test/src/components/Card.jsx":   true,
		"/test/src/utils/helpers.ts":      true,
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("FindFiles() found %d files, expected %d", len(files), len(expectedFiles))
	}

	for _, file := range files {
		if !expectedFiles[file.Path] {
			t.Errorf("FindFiles() found unexpected file: %s", file.Path)
		}
		delete(expectedFiles, file.Path)
	}

	for path := range expectedFiles {
		t.Errorf("FindFiles() missing expected file: %s", path)
	}
}

func TestShouldExcludeDir(t *testing.T) {
	tests := []struct {
		path        string
		excludeDirs []string
		expected    bool
	}{
		{"/src/node_modules", []string{"node_modules"}, true},
		{"/src/components", []string{"node_modules"}, false},
		{"/src/.git", []string{".git", "dist"}, true},
		{"/src/build", []string{"build"}, true},
		{"/src/some/nested/path", []string{"nested"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldExcludeDir(tt.path, tt.excludeDirs)
			if result != tt.expected {
				t.Errorf("shouldExcludeDir(%s, %v) = %v, want %v", tt.path, tt.excludeDirs, result, tt.expected)
			}
		})
	}
}

func TestShouldIncludeFile(t *testing.T) {
	tests := []struct {
		path         string
		filePatterns []string
		expected     bool
	}{
		{"/src/Button.tsx", []string{"*.tsx", "*.ts"}, true},
		{"/src/Button.js", []string{"*.tsx", "*.ts"}, false},
		{"/src/utils.ts", []string{"*.tsx", "*.ts"}, true},
		{"/src/README.md", []string{"*.tsx", "*.ts"}, false},
		{"/src/config.json", []string{"*.json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldIncludeFile(tt.path, tt.filePatterns)
			if result != tt.expected {
				t.Errorf("shouldIncludeFile(%s, %v) = %v, want %v", tt.path, tt.filePatterns, result, tt.expected)
			}
		})
	}
}

func TestShouldExcludeFile(t *testing.T) {
	tests := []struct {
		path         string
		excludeFiles []string
		expected     bool
	}{
		{"/src/Button.test.tsx", []string{"*.test.*"}, true},
		{"/src/Button.tsx", []string{"*.test.*"}, false},
		{"/src/Button.spec.ts", []string{"*.test.*", "*.spec.*"}, true},
		{"/src/Button.stories.tsx", []string{"*.stories.*"}, true},
		{"/src/utils.ts", []string{"*.test.*"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := shouldExcludeFile(tt.path, tt.excludeFiles)
			if result != tt.expected {
				t.Errorf("shouldExcludeFile(%s, %v) = %v, want %v", tt.path, tt.excludeFiles, result, tt.expected)
			}
		})
	}
}