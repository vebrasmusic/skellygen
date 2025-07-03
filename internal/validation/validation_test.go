package validation

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/config"
)

func TestValidateConfig(t *testing.T) {
	originalFs := AppFs
	AppFs = afero.NewMemMapFs()
	defer func() { AppFs = originalFs }()

	AppFs.MkdirAll("/test/src", 0755)

	tests := []struct {
		name        string
		config      *config.Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "/test/src",
					FilePatterns: []string{"*.tsx", "*.ts"},
					ExcludeDirs:  []string{"node_modules"},
					ExcludeFiles: []string{"*.test.*"},
				},
				Output: config.Output{
					NamingPattern: "{component}-skeleton.{ext}",
				},
			},
			expectError: false,
		},
		{
			name: "empty src_dir",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "",
					FilePatterns: []string{"*.tsx"},
				},
				Output: config.Output{
					NamingPattern: "{component}-skeleton.{ext}",
				},
			},
			expectError: true,
			errorMsg:    "src_dir cannot be empty",
		},
		{
			name: "non-existent src_dir",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "/nonexistent",
					FilePatterns: []string{"*.tsx"},
				},
				Output: config.Output{
					NamingPattern: "{component}-skeleton.{ext}",
				},
			},
			expectError: true,
			errorMsg:    "src_dir does not exist: /nonexistent",
		},
		{
			name: "empty file patterns",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "/test/src",
					FilePatterns: []string{},
				},
				Output: config.Output{
					NamingPattern: "{component}-skeleton.{ext}",
				},
			},
			expectError: true,
			errorMsg:    "file_patterns cannot be empty",
		},
		{
			name: "invalid naming pattern - no component",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "/test/src",
					FilePatterns: []string{"*.tsx"},
				},
				Output: config.Output{
					NamingPattern: "skeleton.{ext}",
				},
			},
			expectError: true,
			errorMsg:    "naming_pattern must contain {component} or {name}",
		},
		{
			name: "invalid naming pattern - no extension",
			config: &config.Config{
				Input: config.Input{
					SrcDir:       "/test/src",
					FilePatterns: []string{"*.tsx"},
				},
				Output: config.Output{
					NamingPattern: "{component}-skeleton",
				},
			},
			expectError: true,
			errorMsg:    "naming_pattern must contain {ext}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateConfig(tt.config)
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateConfig() expected error but got none")
					return
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("ValidateConfig() error = %v, want %v", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateConfig() unexpected error = %v", err)
				}
			}
		})
	}
}