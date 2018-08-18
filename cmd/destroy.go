package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/gobuffalo/buffalo-heroku/heroku"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var destroyOptions = struct {
	AppName string
	Branch  string
}{}

// destroyCmd represents the new command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "destroys your application on heroku",
	RunE: func(cmd *cobra.Command, args []string) error {
		app, err := heroku.New(context.Background(), ".")
		if err != nil {
			return errors.WithStack(err)
		}
		c := exec.Command("heroku", "destroy", "--confirm", app.Name.String())
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout
		c.Stdin = os.Stdin
		return c.Run()
	},
}

func init() {
	// auto_cert_mgmt=false
	// addons=heroku-postgresql:hobby-dev
	// collaborators=
	// create_status=undefined
	// git_url=https://git.heroku.com/cryptic-forest-75448.git
	// web_url=https://cryptic-forest-75448.herokuapp.com/
	// repo_size=0 B
	// slug_size=0 B
	// owner=mark@markbates.com
	// region=us
	// dynos={ web: 1 }
	// stack=container
	destroyCmd.Flags().StringVarP(&destroyOptions.AppName, "app-name", "a", "", "the heroku application to delete")
	herokuCmd.AddCommand(destroyCmd)
}
