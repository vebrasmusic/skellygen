package validation

import (
	"errors"
	"os"
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
	if input.SrcDir == "" {
		return errors.New("src_dir cannot be empty")
	}

	if _, err := os.Stat(input.SrcDir); os.IsNotExist(err) {
		return errors.New("src_dir does not exist: " + input.SrcDir)
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
	if output.NamingPattern == "" {
		return errors.New("naming_pattern cannot be empty")
	}

	if !strings.Contains(output.NamingPattern, "{component}") && !strings.Contains(output.NamingPattern, "{name}") {
		return errors.New("naming_pattern must contain {component} or {name}")
	}

	if !strings.Contains(output.NamingPattern, "{ext}") {
		return errors.New("naming_pattern must contain {ext}")
	}

	return nil
}