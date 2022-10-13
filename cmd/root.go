/*
Copyright © 2022 Madhav Jivrajani

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/MadhavJivrajani/cr/gh"
	"github.com/MadhavJivrajani/cr/printer"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

type Options struct {
	PR int
}

var o = Options{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cr",
	Short: "A tool to help reviewers better understand impact of merging a PR.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return processPR(context.Background(), o)
	},
}

func processPR(ctx context.Context, opts Options) error {
	result, err := gh.GetRelevantPullRequests(ctx, opts.PR)
	if err != nil {
		return err
	}

	printer.PrettyDisplayPRs(result)
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().IntVar(&o.PR, "pr", -1, "used to specify the PR number being reviewed")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
