package cmd

import (
	"context"

	"github.com/gobuffalo/buffalo-heroku/genny/heroku"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/gotools"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var generateOptions = struct {
	*heroku.Options
	dryRun bool
}{
	Options: &heroku.Options{},
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "heroku",
	Short: "generates a new heroku",
	RunE: func(cmd *cobra.Command, args []string) error {
		r := genny.WetRunner(context.Background())

		if generateOptions.dryRun {
			r = genny.DryRunner(context.Background())
		}

		g, err := heroku.New(generateOptions.Options)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		g, err = gotools.GoFmt(r.Root)
		if err != nil {
			return errors.WithStack(err)
		}
		r.With(g)

		return r.Run()
	},
}

func init() {
	generateCmd.Flags().BoolVarP(&generateOptions.dryRun, "dry-run", "d", false, "run the generator without creating files or running commands")
	herokuCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(generateCmd)
}
