package validation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/vebrasmusic/skellygen/internal/config"
)

func ValidateConfig(cfg *config.Config) error {
	if err := validateInput(cfg.Input); err != nil {
		return err
	}

	if err := validateOutput(cfg.Output); err != nil {
		return err
	}

	return nil
}

func validateInput(input config.Input) error {
	if input.ReadDir == "" {
		return errors.New("read_dir cannot be empty")
	}

	if _, err := os.Stat(input.ReadDir); os.IsNotExist(err) {
		return errors.New("read_dir does not exist: " + input.ReadDir)
	}

	if len(input.FilePatterns) == 0 {
		return errors.New("file_patterns cannot be empty")
	}

	for _, pattern := range input.FilePatterns {
		if !strings.Contains(pattern, "*") && !strings.Contains(pattern, "?") {
			return errors.New("invalid file pattern: " + pattern)
		}
	}

	return nil
}

func validateOutput(output config.Output) error {
	if output.WriteDir == "" {
		return errors.New("write_dir cannot be empty")
	}

	if output.NamingPattern == "" {
		return errors.New("naming_pattern cannot be empty")
	}

	if !strings.Contains(output.NamingPattern, "{component}") && !strings.Contains(output.NamingPattern, "{name}") {
		return errors.New("naming_pattern must contain {component} or {name}")
	}

	if !strings.Contains(output.NamingPattern, "{ext}") {
		return errors.New("naming_pattern must contain {ext}")
	}

	writeDir := filepath.Dir(output.WriteDir)
	if writeDir != "." {
		if _, err := os.Stat(writeDir); os.IsNotExist(err) {
			return errors.New("parent directory of write_dir does not exist: " + writeDir)
		}
	}

	return nil
}