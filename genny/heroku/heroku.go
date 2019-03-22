package heroku

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gobuffalo/buffalo-heroku/genny/config"
	"github.com/gobuffalo/buffalo-heroku/heroku"
	"github.com/gobuffalo/genny"
	"github.com/gobuffalo/meta"
	"github.com/gobuffalo/x/defaults"
	"github.com/gobuffalo/x/randx"
	"github.com/pkg/errors"
)

func New(opts *Options) (*genny.Group, error) {
	gg := &genny.Group{}

	conf, err := config.New(&config.Options{})
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(conf)

	g, err := build(opts)
	if err != nil {
		return gg, errors.WithStack(err)
	}
	gg.Add(g)
	return gg, nil
}

func build(opts *Options) (*genny.Generator, error) {
	g := genny.New()

	if !isValid(opts.DynoLevel, heroku.DynoLevels) {
		return g, errors.Errorf("%s is not a valid dyno level", opts.DynoLevel)
	}

	if opts.Auth {
		g.Command(exec.Command("heroku", "login"))
		g.Command(exec.Command("heroku", "container:login"))
	}
	g.RunFn(installHerokuCLI)

	g.Command(exec.Command("heroku", "create", "--manifest", opts.AppName, "--region", opts.Region))
	g.Command(exec.Command("heroku", "stack:set", "container"))

	configs := map[string]string{
		"GO_ENV":         defaults.String(opts.Environment, "production"),
		"SESSION_SECRET": randx.String(100),
	}

	for k, v := range configs {
		g.Command(exec.Command("heroku", "config:set", fmt.Sprintf("%s=%s", k, v)))
	}

	for _, a := range opts.Addons {
		if len(a.Level) == 0 {
			continue
		}
		if !isValid(a.Level, a.Available) {
			return g, errors.Errorf("%s is not a valid level for %s", a.Level, a.Name)
		}
		func(a heroku.Addon) {
			g.RunFn(func(r *genny.Runner) error {
				r.Logger.Info("addon:started", a)
				if err := r.Exec(exec.Command("heroku", "addons:create", fmt.Sprintf("%s:%s", a.Name, a.Level))); err != nil {
					return errors.WithStack(err)
				}
				r.Logger.Info("addon:finished", a)
				if a.Name == "sendgrid" {
					return setupSendgrid(r)
				}
				return nil
			})
		}(a)
	}
	g.RunFn(func(r *genny.Runner) error {
		var files []string
		for _, f := range r.Disk.Files() {
			files = append(files, f.Name())
		}
		if len(files) > 0 {
			cmd := exec.Command("git", "add")
			cmd.Args = append(cmd.Args, files...)
			if err := r.Exec(cmd); err != nil {
				return errors.WithStack(err)
			}

			cmd = exec.Command("git", "commit", "-m", "heroku config")
			cmd.Args = append(cmd.Args, files...)
			bb := &bytes.Buffer{}
			cmd.Stderr = bb
			cmd.Stdout = bb
			if err := r.Exec(cmd); err != nil {
				if !strings.Contains(bb.String(), "nothing to commit, working tree clean") {
					return errors.WithStack(err)
				}
			}
		}
		return nil
	})

	app := meta.New(".")
	if app.WithPop {
		g.RunFn(func(r *genny.Runner) error {
			gk := filepath.Join(app.Root, "migrations", ".git-keep")
			if _, err := os.Stat(gk); err == nil {
				return nil
			}
			f := genny.NewFile(filepath.Join(app.Root, "migrations", ".git-keep"), strings.NewReader(""))
			if err := r.File(f); err != nil {
				return errors.WithStack(err)
			}
			cmd := exec.Command("git", "add", f.Name())
			if err := r.Exec(cmd); err != nil {
				return errors.WithStack(err)
			}

			bb := &bytes.Buffer{}
			cmd = exec.Command("git", "commit", "-m", "migrations folder stub")
			cmd.Stderr = bb
			cmd.Stdout = bb
			if err := r.Exec(cmd); err != nil {
				if !strings.Contains(bb.String(), "nothing to commit, working tree clean") {
					return errors.WithStack(err)
				}
			}
			return nil
		})
	}

	// git rev-parse --abbrev-ref HEAD
	g.RunFn(func(r *genny.Runner) error {
		cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		b, err := cmd.CombinedOutput()
		if err != nil {
			return errors.WithStack(err)
		}
		branch := strings.TrimSpace(string(b))
		if branch == "" {
			branch = "master"
		}
		cmd = exec.Command("git", "push", "heroku", branch+":master")
		return r.Exec(cmd)
	})
	g.Command(exec.Command("heroku", "dyno:type", opts.DynoLevel))
	g.Command(exec.Command("heroku", "config"))
	return g, nil
}

func installHerokuCLI(r *genny.Runner) error {
	if _, err := exec.LookPath("heroku"); err != nil {
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("brew"); err == nil {
				c := exec.Command("brew", "install", "heroku")
				return r.Exec(c)
			}
		}
		return errors.New("heroku cli is not installed. https://devcenter.heroku.com/articles/heroku-cli")
	}
	return nil
}

func setupSendgrid(r *genny.Runner) error {
	if err := r.Exec(exec.Command("heroku", "config:set", "SMTP_HOST=smtp.sendgrid.net", "SMTP_PORT=465")); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("heroku", "config:get", "SENDGRID_USERNAME")
	b, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}
	user := strings.TrimSpace(string(b))
	if err := r.Exec(exec.Command("heroku", "config:set", fmt.Sprintf("SMTP_USER=%s", user))); err != nil {
		return errors.WithStack(err)
	}

	cmd = exec.Command("heroku", "config:get", "SENDGRID_PASSWORD")
	b, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Print(string(b))
		return errors.WithStack(err)
	}
	pass := strings.TrimSpace(string(b))
	if err := r.Exec(exec.Command("heroku", "config:set", fmt.Sprintf("SMTP_PASSWORD=%s", pass))); err != nil {
		return errors.WithStack(err)
	}

	return r.Exec(exec.Command("heroku", "config:set", "SMTP_PORT=25"))
}

func isValid(s string, a []string) bool {
	for _, x := range a {
		if s == x {
			return true
		}
	}
	return false
}
