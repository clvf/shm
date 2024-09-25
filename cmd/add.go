package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/urfave/cli/v2"
)

func add(cCtx *cli.Context) error {
	file, err := getFirstPos(cCtx)
	check(err)

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return cli.Exit("EDITOR environment variable must be set!", 1)
	}

	err = os.MkdirAll(filepath.Dir(file), 0750)
	check(err)

	binary, err := exec.LookPath(editor)
	check(err)

	env := os.Environ()
	args := []string{editor, file}

	err = syscall.Exec(binary, args, env)
	check(err)

	return nil
}
