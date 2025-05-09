/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hacker65536/aft-cli/pkg/myaws"
	"github.com/spf13/cobra"
)

// pipelineStatusCmd represents the pipelineStatus command
var pipelineStatusCmd = &cobra.Command{
	Use:   "pipelineStatus",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Aliases: []string{"pipeline-status"},
	Run: func(cmd *cobra.Command, args []string) {
		myaws := myaws.New()
		myaws.AftPipelineStatus()
	},
}

func init() {
	rootCmd.AddCommand(pipelineStatusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pipelineStatusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pipelineStatusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
