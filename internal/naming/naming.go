package naming

import (
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/config"
	"github.com/vebrasmusic/skellygen/internal/discovery"
)

var AppFs afero.Fs = afero.NewOsFs()

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
	return AppFs.MkdirAll(dir, 0755)
}