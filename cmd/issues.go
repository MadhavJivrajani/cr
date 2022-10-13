/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/MadhavJivrajani/cr/gh"
	"github.com/MadhavJivrajani/cr/printer"
	"github.com/spf13/cobra"
)

// issuesCmd represents the issues command
var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "Command to get stats on triage/accepted and kind/bug issues.",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := gh.HowManyIssuesWithLabels(context.Background(), "triage/accepted", "kind/bug")
		if err != nil {
			return err
		}
		printer.PrettyDisplayIssues(res)
		return err
	},
}

func init() {
	rootCmd.AddCommand(issuesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
