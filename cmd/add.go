package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
)

func add(cCtx *cli.Context) error {
	if cCtx.Args().Present() != true {
		return cli.Exit("At least one positional argument is required!", 1)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return cli.Exit("EDITOR environment variable must be set!", 1)
	}

	file, err := filepath.Abs(filepath.Join(_REPO, cCtx.Args().First()))
	check(err)

	dirname := filepath.Dir(file)
	err = os.MkdirAll(dirname, 0750)
	check(err)

	binary, err := exec.LookPath(editor)
	check(err)

	env := os.Environ()
	args := []string{editor, file}

	err = syscall.Exec(binary, args, env)
	check(err)

	return nil
}
