package generation

import (
	"errors"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
	"github.com/vebrasmusic/skellygen/internal/naming"
	"github.com/vebrasmusic/skellygen/internal/validation"
)

func getDirectories() (*config.Config, error) {
	var config config.Config

	yamlFile, err := os.ReadFile("skelly.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("skelly.yaml not found. Please run 'skelly init' first to create a configuration file")
		}
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

	err = validation.ValidateConfig(config)
	if err != nil {
		return err
	}

	files, err := discovery.FindFiles(config)
	if err != nil {
		return err
	}

	for _, file := range files {
		outputPath := naming.GenerateOutputPath(file, config)

		err := naming.EnsureOutputDir(outputPath)
		if err != nil {
			return err
		}

		err = generateSkeleton(file, outputPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateSkeleton(file discovery.FileInfo, outputPath string) error {
	content, err := os.ReadFile(file.Path)
	if err != nil {
		return err
	}

	skeletonContent := string(content)

	return os.WriteFile(outputPath, []byte(skeletonContent), 0644)
}
