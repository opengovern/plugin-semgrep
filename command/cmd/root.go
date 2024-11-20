/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "OpenGovernance aws describer manual",
	RunE: func(cmd *cobra.Command, args []string) error {
		var items []string
		items = append(items, "describer")

		return describerCmd.Help()

	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(describerCmd)
	rootCmd.AddCommand(getDescriberCmd)
}
