package cmd

import (
	"context"

	"github.com/gobuffalo/buffalo-heroku/genny/config"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/gogen"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var configOptions = struct {
	*config.Options
	dryRun bool
}{
	Options: &config.Options{},
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "configs configuration files for heroku",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if configOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		g, err := config.New(configOptions.Options)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gogen.Fmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	configCmd.Flags().BoolVarP(&configOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	herokuCmd.AddCommand(configCmd)
}
