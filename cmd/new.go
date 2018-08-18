package cmd

import (
	"context"
	"strings"

	"github.com/gobuffalo/buffalo-heroku/genny/heroku"
	her "github.com/gobuffalo/buffalo-heroku/heroku"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var newOptions = struct {
	*heroku.Options
	dryRun bool
}{
	Options: &heroku.Options{},
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "creates a new heroku application",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if newOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		opts := newOptions.Options
		opts.Addons = her.DefaultAddons
		gg, err := heroku.New(opts)
		if err != nil {
			return errors.WithStack(err)
		}
		gg.With(r)

		g, err := gotools.GoFmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	newCmd.Flags().StringVarP(&newOptions.DynoLevel, "dyno-level", "l", "free", strings.Join(her.DynoLevels, ", "))
	newCmd.Flags().StringVarP(&newOptions.AppName, "app-name", "a", "", "the name of the heroku app to deploy")
	newCmd.Flags().StringVarP(&newOptions.Environment, "environment", "e", "production", "the environment to run the application in")
	newCmd.Flags().BoolVar(&newOptions.Auth, "auth", false, "log into heroku from the cli")
	newCmd.Flags().BoolVarP(&newOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	herokuCmd.AddCommand(newCmd)
}
