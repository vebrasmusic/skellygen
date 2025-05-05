/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vebrasmusic/skellygen/internal/config"
)

var (
	ReadDir  string
	WriteDir string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates a config file w/ defaults in the project if needed.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := config.RunInit(ReadDir, WriteDir)
		if err != nil {
			return nil
		}

		return nil
	},
}

func init() {
	initCmd.PersistentFlags().StringVarP(&ReadDir, "read-dir", "r", "", "Directory for skellygen to look for files to create skeletons for.")
	initCmd.PersistentFlags().StringVarP(&WriteDir, "write-dir", "w", "", "Directory for skellygen to write skeletons in.")
	initCmd.MarkPersistentFlagRequired("read-dir")
	initCmd.MarkPersistentFlagRequired("write-dir")
	rootCmd.AddCommand(initCmd)
}
