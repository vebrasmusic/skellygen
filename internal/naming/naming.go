package naming

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
)

func GenerateOutputPath(file discovery.FileInfo, cfg *config.Config) string {
	pattern := cfg.Output.NamingPattern
	
	pattern = strings.ReplaceAll(pattern, "{component}", file.Name)
	pattern = strings.ReplaceAll(pattern, "{name}", file.Name)
	pattern = strings.ReplaceAll(pattern, "{ext}", file.Extension)
	
	dir := filepath.Dir(file.Path)
	return filepath.Join(dir, pattern)
}

func EnsureOutputDir(outputPath string) error {
	dir := filepath.Dir(outputPath)
	return os.MkdirAll(dir, 0755)
}