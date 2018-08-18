package cmd

import (
	"fmt"

	"github.com/gobuffalo/buffalo-heroku/heroku"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of heroku",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("heroku", heroku.Version)
		return nil
	},
}

func init() {
	herokuCmd.AddCommand(versionCmd)
}
