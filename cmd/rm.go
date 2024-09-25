package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
)

func rm(cCtx *cli.Context) error {
	file, err := getFirstPos(cCtx)
	check(err)

	err = os.Remove(file)
	check(err)

	return nil
}
