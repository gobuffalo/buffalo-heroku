package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "v0.0.0-dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "buffalo-heroku",
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Usage()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
