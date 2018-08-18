package heroku

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
	"time"

	bmeta "github.com/gobuffalo/buffalo/meta"
	"github.com/markbates/inflect"
)

type App struct {
	context.Context
	Name        inflect.Name
	Addons      Addons
	Environment string
	DynoLevel   string
	Auth        bool
	Existing    bool
	Buffalo     bmeta.App
}

func New(ctx context.Context, root string) (*App, error) {
	b := bmeta.New(root)
	app := &App{
		Context:     ctx,
		Buffalo:     b,
		Environment: "production",
		DynoLevel:   DynoLevels[0],
		Addons:      DefaultAddons,
	}

	if app.Name.String() == "" {
		n, e := herokuNameFromGit(app)
		app.Name = inflect.Name(n)
		app.Existing = e
	}

	return app, nil
}

var gitHerokuRX = regexp.MustCompile("heroku.com/(.+).git")

func herokuNameFromGit(ctx context.Context) (string, bool) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, "git", "remote", "-v")
	out, err := c.CombinedOutput()
	if err != nil {
		return "", false
	}

	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "heroku") {
			continue
		}

		names := gitHerokuRX.FindStringSubmatch(line)
		if len(names) > 1 {
			return names[1], true
		}
	}

	return "", false
}
