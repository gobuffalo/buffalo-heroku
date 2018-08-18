package heroku

import (
	"github.com/gobuffalo/buffalo/meta"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"

	br "github.com/gobuffalo/buffalo/runtime"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := g.Box(packr.NewBox("../heroku/templates")); err != nil {
		return g, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)

	app := meta.New(".")
	ctx.Set("app", app)

	ctx.Set("version", br.Version)
	g.Transformer(plushgen.Transformer(ctx))
	return g, nil
}
