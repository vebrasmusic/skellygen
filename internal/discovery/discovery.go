package discovery

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/afero"
	"github.com/vebrasmusic/skellygen/internal/config"
)

var AppFs afero.Fs = afero.NewOsFs()

type FileInfo struct {
	Path         string
	RelativePath string
	Name         string
	Extension    string
}

func FindFiles(cfg *config.Config) ([]FileInfo, error) {
	var files []FileInfo

	err := afero.Walk(AppFs, cfg.Input.SrcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if shouldExcludeDir(path, cfg.Input.ExcludeDirs) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldIncludeFile(path, cfg.Input.FilePatterns) && !shouldExcludeFile(path, cfg.Input.ExcludeFiles) {
			relativePath, err := filepath.Rel(cfg.Input.SrcDir, path)
			if err != nil {
				return err
			}

			files = append(files, FileInfo{
				Path:         path,
				RelativePath: relativePath,
				Name:         strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
				Extension:    strings.TrimPrefix(filepath.Ext(info.Name()), "."),
			})
		}

		return nil
	})

	return files, err
}

func shouldExcludeDir(path string, excludeDirs []string) bool {
	dirName := filepath.Base(path)
	return slices.Contains(excludeDirs, dirName)
}

func shouldIncludeFile(path string, filePatterns []string) bool {
	fileName := filepath.Base(path)
	for _, pattern := range filePatterns {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}
	return false
}

func shouldExcludeFile(path string, excludeFiles []string) bool {
	fileName := filepath.Base(path)
	for _, pattern := range excludeFiles {
		if matched, _ := filepath.Match(pattern, fileName); matched {
			return true
		}
	}
	return false
}

