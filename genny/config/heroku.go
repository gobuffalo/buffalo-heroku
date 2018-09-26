package config

import (
	"bytes"
	"strings"

	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/genny/movinglater/plushgen"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/pkg/errors"

	br "github.com/gobuffalo/buffalo/runtime"
)

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	if err := g.Box(packr.NewBox("../config/templates")); err != nil {
		return g, errors.WithStack(err)
	}
	ctx := plush.NewContext()
	ctx.Set("opts", opts)
	ctx.Set("app", opts.App)
	ctx.Set("version", br.Version)

	g.Transformer(plushgen.Transformer(ctx))

	g.RunFn(func(r *genny.Runner) error {
		f, err := r.FindFile("Dockerfile")
		bb := &bytes.Buffer{}
		if err != nil {
			return errors.WithStack(err)
		}
		for _, line := range strings.Split(f.String(), "\n") {
			if strings.HasPrefix(line, "FROM alpine") {
				line += "\nRUN apk add --no-cache curl"
			}
			bb.WriteString(line + "\n")
		}
		return r.File(genny.NewFile(f.Name(), bb))
	})

	return g, nil
}
