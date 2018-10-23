package config

import "github.com/gobuffalo/meta"

type Options struct {
	App meta.App
}

func (opts *Options) Validate() error {
	if opts.App.IsZero() {
		opts.App = meta.New(".")
	}
	return nil
}
