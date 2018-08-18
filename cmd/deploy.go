package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var deployOptions = struct {
	AppName string
	Branch  string
}{}

// deployCmd represents the new command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("git", "push", deployOptions.AppName, fmt.Sprintf("%s:master", deployOptions.Branch))
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		return c.Run()
	},
}

func init() {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	o, err := cmd.CombinedOutput()
	if err != nil {
		o = []byte("master")
	}
	branch := string(o)
	branch = strings.TrimSpace(branch)
	deployCmd.Flags().StringVarP(&deployOptions.AppName, "app-name", "a", "heroku", "the git end point to push to")
	deployCmd.Flags().StringVarP(&deployOptions.Branch, "branch", "b", branch, "the name of the branch to depoly")
	herokuCmd.AddCommand(deployCmd)
}
