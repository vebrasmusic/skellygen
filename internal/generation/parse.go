package generation

import (
	"os"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/goccy/go-yaml"
	"github.com/vebrasmusic/skellygen/internal/config"
)

func getDirectories() (*config.Config, error) {
	var config config.Config

	yamlFile, err := os.ReadFile("skelly.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func ParseInputFile() error {
	config, err := getDirectories()
	if err != nil {
		return err
	}

	path := filepath.Join(config.ReadDir, "test.tsx")

	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{path},
		Outfile:     filepath.Join(config.WriteDir, "out.js"),
		Bundle:      true,
		Write:       true,
		LogLevel:    esbuild.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
	return nil
}
