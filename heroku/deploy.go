package heroku

import (
	"fmt"
	"os"
	"os/exec"
)

func Deploy(name string, branch string) error {
	c := exec.Command("git", "push", name, fmt.Sprintf("%s:master", branch))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	return c.Run()
}
