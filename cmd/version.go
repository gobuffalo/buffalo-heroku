package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "current version of heroku",
	RunE: func(cmd *cobra.Command, args []string) error {
		o, _ := exec.Command("heroku", "--version").Output()
		fmt.Println("heroku", string(o))
		return nil
	},
}

func init() {
	herokuCmd.AddCommand(versionCmd)
}
