package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

const pathSeparator = string(filepath.Separator)

func visitor(pattern string) func(string, fs.DirEntry, error) error {

	// NOTE: the following check is bogus since *,?,[ etc. can be escaped and
	//		 we don't detect that.
	patternIsGlob := strings.ContainsAny(pattern, "*?[")

	searchPattern := pattern
	if !patternIsGlob {
		searchPattern = fmt.Sprintf("*%s*", pattern)
	}

	fmt.Println("Search Terms:", searchPattern)

	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// _REPO itself doesn't need to be checked
		if path == _REPO {
			return nil
		}

		rel, err := filepath.Rel(_REPO, path)
		if err != nil {
			return err
		}

		matchFound := false
		parts := strings.Split(rel, pathSeparator)
		for _, p := range parts {

			matched, err := filepath.Match(searchPattern, p)
			if err != nil {
				return err
			}

			if matched {
				matchFound = true
			}

		}

		if matchFound {
			fmt.Println(rel)
		}

		return nil

	}
}

func search(cCtx *cli.Context) error {
	if cCtx.Args().Present() != true {
		return cli.Exit("At least one positional argument is required!", 1)
	}

	err := filepath.WalkDir(_REPO, visitor(cCtx.Args().First()))
	check(err)

	return nil
}
