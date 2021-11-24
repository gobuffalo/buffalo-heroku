package config

import (
	"bytes"
	"embed"
	"io/fs"
	"strings"

	"github.com/gobuffalo/genny/v2"
	"github.com/gobuffalo/genny/v2/plushgen"
	"github.com/gobuffalo/plush/v4"
	"github.com/pkg/errors"

	br "github.com/gobuffalo/buffalo/runtime"
)

//go:embed templates
var templates embed.FS

func New(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if err := opts.Validate(); err != nil {
		return g, errors.WithStack(err)
	}

	sub, err := fs.Sub(templates, "templates")
	if err != nil {
		return g, errors.WithStack(err)
	}
	if err := g.FS(sub); err != nil {
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
