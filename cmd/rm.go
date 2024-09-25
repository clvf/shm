package cmd

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func rm(cCtx *cli.Context) error {
	if cCtx.Args().Present() != true {
		return cli.Exit("At least one positional argument is required!", 1)
	}

	file, err := filepath.Abs(filepath.Join(_REPO, cCtx.Args().First()))
	check(err)

	err = os.Remove(file)
	check(err)

	return nil
}
