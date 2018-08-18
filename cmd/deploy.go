package cmd

import (
	"os/exec"
	"strings"

	"github.com/gobuffalo/buffalo-heroku/heroku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deployOptions = struct {
	Name   string
	Branch string
}{}

// deployCmd represents the new command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "deploys your application",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := heroku.Deploy(deployOptions.Name, deployOptions.Branch); err != nil {
			return errors.WithStack(err)
		}
		c := exec.Command("heroku", "open")
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
	deployCmd.Flags().StringVarP(&deployOptions.Name, "name", "n", "heroku", "the git end point to push to")
	deployCmd.Flags().StringVarP(&deployOptions.Branch, "branch", "b", branch, "the name of the branch to depoly")
	herokuCmd.AddCommand(deployCmd)
}
