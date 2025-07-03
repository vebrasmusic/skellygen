package naming

import (
	"testing"

	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
)

func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		name     string
		file     discovery.FileInfo
		config   *config.Config
		expected string
	}{
		{
			name: "basic component skeleton",
			file: discovery.FileInfo{
				Path:         "/src/components/Button.tsx",
				RelativePath: "components/Button.tsx",
				Name:         "Button",
				Extension:    "tsx",
			},
			config: &config.Config{
				Output: config.Output{
					NamingPattern: "{component}-skeleton.{ext}",
				},
			},
			expected: "/src/components/Button-skeleton.tsx",
		},
		{
			name: "custom naming pattern",
			file: discovery.FileInfo{
				Path:         "/src/Card.js",
				RelativePath: "Card.js",
				Name:         "Card",
				Extension:    "js",
			},
			config: &config.Config{
				Output: config.Output{
					NamingPattern: "{name}-loading.{ext}",
				},
			},
			expected: "/src/Card-loading.js",
		},
		{
			name: "nested component",
			file: discovery.FileInfo{
				Path:         "/src/ui/forms/Input.tsx",
				RelativePath: "ui/forms/Input.tsx",
				Name:         "Input",
				Extension:    "tsx",
			},
			config: &config.Config{
				Output: config.Output{
					NamingPattern: "{component}.skeleton.{ext}",
				},
			},
			expected: "/src/ui/forms/Input.skeleton.tsx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateOutputPath(tt.file, tt.config)
			if result != tt.expected {
				t.Errorf("GenerateOutputPath() = %v, want %v", result, tt.expected)
			}
		})
	}
}