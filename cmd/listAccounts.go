/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hacker65536/aft-cli/pkg/myaws"
	"github.com/spf13/cobra"
)

// listAccountsCmd represents the listAccounts command
var listAccountsCmd = &cobra.Command{
	Use:   "listAccounts",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"list-accounts"},
	Run: func(cmd *cobra.Command, args []string) {
		myaws := myaws.New()
		myaws.ListAccounts()
		//myaws.ListCodePipelines()
		//myaws.AftPipelineStatus()
	},
}

func init() {
	rootCmd.AddCommand(listAccountsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listAccountsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listAccountsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
