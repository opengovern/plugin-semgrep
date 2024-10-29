/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aws",
	Short: "OpenGovernance aws describer manual",
	RunE: func(cmd *cobra.Command, args []string) error {
		var items []string
		items = append(items, "describer")
		items = append(items, "getDescriber")
		prompt := promptui.Select{
			Label: "Please select the types of describer",
			Items: items,
		}
		_, result, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("[workspaces] : %v", err)
		}
		typeDescriber := result
		if typeDescriber == "describer" {
			return describerCmd.Help()
		} else {
			return getDescriberCmd.Help()
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(getDescriberCmd)
	rootCmd.AddCommand(describerCmd)
}
