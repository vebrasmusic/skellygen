package config

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/vebrasmusic/skellygen/internal/utils"
)

type Config struct {
	ReadDir  string
	WriteDir string
}

func RunInit(WriteDir string, ReadDir string) error {
	configExists, err := utils.CheckForConfig()
	if err != nil {
		return err
	}
	if configExists {
		panic("Config already exists, stopping.")
	}

	data := Config{
		ReadDir,
		WriteDir,
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
