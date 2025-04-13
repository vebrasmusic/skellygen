/*
Copyright © 2025 Andres Duvvuri vebrasmusic@gmail.com
*/
package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vebrasmusic/skellygen/internal/generation"
)

var (
	OutputPath string
	InputFile  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "skellygen",
	Short: "Generate pixel-perfect loading skeletons for your components.",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		outputPath, err := filepath.Abs(OutputPath)
		if err != nil {
			panic("Some error")
		}
		log.Println("output file set to: ", outputPath)
		inputPath, err := filepath.Abs(InputFile)
		if err != nil {
			panic("Some error")
		}
		log.Println("input file set to: ", inputPath)

		generation.ParseInputFile(inputPath, outputPath)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&OutputPath, "output", "o", "skeletons.tsx", "Sets output directory for skeleton generation.")
	rootCmd.PersistentFlags().StringVarP(&InputFile, "input", "i", "", "Sets input file to generate skeletons against.")
	rootCmd.MarkPersistentFlagRequired("input")
}
