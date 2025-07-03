/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vebrasmusic/skellygen/internal/config"
)

var (
	SrcDir         string
	FilePatterns   string
	ExcludeDirs    string
	ExcludeFiles   string
	NamingPattern  string
	PreserveStruct bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates a config file w/ defaults in the project if needed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.RunInit(SrcDir, FilePatterns, ExcludeDirs, ExcludeFiles, NamingPattern, PreserveStruct)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&SrcDir, "src-dir", "s", "./src", "Directory for skellygen to look for files to create skeletons for.")
	initCmd.PersistentFlags().StringVar(&FilePatterns, "patterns", "", "Comma-separated file patterns to include (e.g., '*.tsx,*.ts')")
	initCmd.PersistentFlags().StringVar(&ExcludeDirs, "exclude-dirs", "", "Comma-separated directories to exclude (e.g., 'node_modules,dist')")
	initCmd.PersistentFlags().StringVar(&ExcludeFiles, "exclude-files", "", "Comma-separated file patterns to exclude (e.g., '*.test.*,*.spec.*')")
	initCmd.PersistentFlags().StringVar(&NamingPattern, "naming-pattern", "", "Output file naming pattern (e.g., '{component}-skeleton.{ext}')")
	initCmd.PersistentFlags().BoolVar(&PreserveStruct, "preserve-structure", true, "Preserve directory structure in output")
	rootCmd.AddCommand(initCmd)
}
