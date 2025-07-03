package config

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/vebrasmusic/skellygen/internal/utils"
)

type Project struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Input struct {
	SrcDir       string   `yaml:"src_dir"`
	FilePatterns []string `yaml:"file_patterns"`
	ExcludeDirs  []string `yaml:"exclude_dirs"`
	ExcludeFiles []string `yaml:"exclude_files"`
}

type Output struct {
	NamingPattern     string `yaml:"naming_pattern"`
	PreserveStructure bool   `yaml:"preserve_structure"`
}

type Config struct {
	Project Project `yaml:"project"`
	Input   Input   `yaml:"input"`
	Output  Output  `yaml:"output"`
}

func RunInit(SrcDir string, FilePatterns string, ExcludeDirs string, ExcludeFiles string, NamingPattern string, PreserveStruct bool) error {
	configExists, err := utils.CheckForConfig()
	if err != nil {
		return err
	}
	if configExists {
		panic("Config already exists, stopping.")
	}

	filePatterns := []string{"*.tsx", "*.ts", "*.jsx", "*.js"}
	if FilePatterns != "" {
		filePatterns = strings.Split(FilePatterns, ",")
	}

	excludeDirs := []string{"node_modules", ".git", "dist", "build", ".next"}
	if ExcludeDirs != "" {
		excludeDirs = strings.Split(ExcludeDirs, ",")
	}

	excludeFiles := []string{"*.test.*", "*.spec.*", "*.stories.*"}
	if ExcludeFiles != "" {
		excludeFiles = strings.Split(ExcludeFiles, ",")
	}

	namingPattern := "{component}-skeleton.{ext}"
	if NamingPattern != "" {
		namingPattern = NamingPattern
	}

	data := Config{
		Project: Project{
			Name:    "my-project",
			Version: "1.0.0",
		},
		Input: Input{
			SrcDir:       SrcDir,
			FilePatterns: filePatterns,
			ExcludeDirs:  excludeDirs,
			ExcludeFiles: excludeFiles,
		},
		Output: Output{
			NamingPattern:     namingPattern,
			PreserveStructure: PreserveStruct,
		},
	}

	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	file, err := os.Create("skelly.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString(string(bytes))

	return nil
}
