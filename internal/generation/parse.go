package generation

import (
	"log"
	"os"
	"path/filepath"

	esbuild "github.com/evanw/esbuild/pkg/api"
)

func ParseInputFile(inputFile string, outputDir string) {
	data, err := os.ReadFile(inputFile)
	Check(err)
	log.Println("Rec'vd data: ", string(data))

	result := esbuild.Build(esbuild.BuildOptions{
		EntryPoints: []string{inputFile},
		Outfile:     filepath.Join(outputDir, "out.js"),
		Bundle:      true,
		Write:       true,
		LogLevel:    esbuild.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}
}
