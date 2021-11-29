package cmd

import (
	"github.com/spf13/cobra"
)

// herokuCmd represents the heroku command
var herokuCmd = &cobra.Command{
	Use:   "heroku",
	Short: "helps with heroku setup and deployment for buffalo applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(herokuCmd)
}
