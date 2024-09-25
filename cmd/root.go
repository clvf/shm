package cmd

import (
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

const _BASEDIR string = ".shm" // the basename of the repository (see shm_repo)
const _COLOR_RESET = "\033[0m"
const _COLOR_BLUE = "\033[34m"
const _DEFAULT_WIDTH = 120
const _DEFAULT_STYLE = "tokyo-night"

var _APP *cli.App
var _REPO string // the directory where the snippets (files) are stored
var _MONOCHROME bool

func check(err interface{}) {
	if err != nil {
		log.Fatalln(err)
	}
}

// Return the default path to place the files to: $HOME/<shm_basedir>
func getRepoPath() string {
	var repoPath string

	if home := os.Getenv("HOME"); home != "" {
		return filepath.Join(home, _BASEDIR)
	}

	// HOME couldn't be found.
	usr, err := user.Current()
	check(err)

	repoPath = filepath.Join("/var/tmp", usr.Name, _BASEDIR)
	log.Println("Cannot find your HOME directory. Fallback dir: %q", repoPath)

	return repoPath
}

// Ensure that there's a positional parameter and return it with it's absolute path.
func getFirstPos(cCtx *cli.Context) (string, error) {
	if cCtx.Args().Present() != true {
		return "", cli.Exit("At least one positional argument is required!", 1)
	}

	abspath, err := filepath.Abs(filepath.Join(_REPO, cCtx.Args().First()))
	check(err)

	return abspath, nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_APP = &cli.App{
		Name:                   "shm",
		Usage:                  "SHow Me - manage code snippets and command recipes as a directory tee.",
		Action:                 show,
		ArgsUsage:              "[FILE]",
		Description:            "Show the snippet designated by FILE in the repository or if no arguments given print the repository's content",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "repository",
				Aliases:     []string{"r"},
				Destination: &_REPO,
				Value:       getRepoPath(),
				Usage:       "The directory (`REPO`) where the snippets (files) are stored",
				EnvVars:     []string{"SHM_REPO"},
			},
			&cli.BoolFlag{
				Name:        "monochrome",
				Aliases:     []string{"m"},
				Destination: &_MONOCHROME,
				Usage:       "Print output without colors",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "Print the contents of the repository",
				Action: func(cCtx *cli.Context) error {
					printRepo(_REPO)
					return nil
				},
			},
			{
				Name:        "add",
				Aliases:     []string{"a", "e", "edit"},
				Usage:       "Add a new recipe",
				Description: "Open $EDITOR and add a new recipe in form of FILE under repository.",
				ArgsUsage:   "FILE",
				Action:      add,
			},
			{
				Name:        "search",
				Aliases:     []string{"s", "f", "find"},
				Usage:       "Search for a recipe",
				Description: "Print files in the repository directory that match the glob in params.",
				ArgsUsage:   "GLOB",
				Action:      search},
			{
				Name:        "rm",
				Usage:       "Delete a recipe",
				Description: "Delete the file from the repository directory",
				ArgsUsage:   "FILE",
				Action:      rm,
			},
		},
	}
}

func Run() {
	if err := _APP.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
