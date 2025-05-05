/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vebrasmusic/skellygen/internal/generation"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generates the skeletons based on config in yaml.",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := generation.ParseInputFile()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
}
