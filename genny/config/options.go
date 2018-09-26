package config

import "github.com/gobuffalo/buffalo/meta"

type Options struct {
	App meta.App
}

func (opts *Options) Validate() error {
	if (opts.App == meta.App{}) {
		opts.App = meta.New(".")
	}
	return nil
}
