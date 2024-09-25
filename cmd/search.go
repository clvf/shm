package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func search(cCtx *cli.Context) error {
	if cCtx.Args().Present() != true {
		return cli.Exit("At least one positional argument is required!", 1)
	}
	pattern := cCtx.Args().First()

	if filepath.IsAbs(pattern) == false {
		pattern = filepath.Join(_REPO, pattern)
	}

	matches, err := filepath.Glob(pattern)
	check(err)

	fmt.Println("%v", matches)

	return nil
}
