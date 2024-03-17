/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run rogue-like poker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(); err != nil {
			return err
		}
		return nil
	},
}

func run() error {
	fmt.Println("run called")
	return nil
}

func init() {
	rootCmd.AddCommand(runCmd)
}
