package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/glamour"
	"github.com/urfave/cli/v2"
)

const _BASEDIR string = ".shm" // the basename of the repository (see shm_repo)
const COLOR_RESET = "\033[0m"
const COLOR_BLUE = "\033[34m"
const _DEFAULTWIDTH = 120

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

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_APP = &cli.App{
		Name:                   "shm",
		Usage:                  "SHow Me - manage code snippets and command recipes as a directory tee.",
		Action:                 show_or_list,
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
					printRepo()
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

func show_or_list(cCtx *cli.Context) error {
	if !cCtx.Args().Present() {
		printRepo()
		return nil
	}

	file, err := filepath.Abs(filepath.Join(_REPO, cCtx.Args().First()))
	check(err)

	content, err := os.ReadFile(file)
	check(err)

	if !_MONOCHROME {
		style := "tokyo-night"
		if override := os.Getenv("GLAMOUR_STYLE"); override != "" {
			style = override
		}

		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle(style),
			glamour.WithWordWrap(_DEFAULTWIDTH),
		)
		check(err)

		out, err := r.Render(string(content))
		check(err)

		fmt.Println(out)
	} else {
		os.Stdout.Write(content)
	}

	return nil
}

func printRepo() {
	fmt.Println(_REPO)

	var sb strings.Builder
	walkDir(&sb, _REPO, "")
	fmt.Print(sb.String())
}

func walkDir(sb_ptr *strings.Builder, root string, linePrefix string) {
	files, err := os.ReadDir(root)
	check(err)

	fileNum := len(files)
	for i, file := range files {
		filePrefix := "\u2514\u2500 " // "|- "

		if i < fileNum-1 {
			filePrefix = "\u251c\u2500 " // "`- "
		}

		if file.IsDir() && !_MONOCHROME {
			fmt.Fprintf(sb_ptr, "%v%v%v%v%v\n", linePrefix, filePrefix, COLOR_BLUE, file.Name(), COLOR_RESET)
		} else {
			fmt.Fprintf(sb_ptr, "%v%v%v\n", linePrefix, filePrefix, file.Name())
		}

		if file.IsDir() {
			nextLinePrefix := linePrefix + "   "
			if i < fileNum-1 {
				nextLinePrefix = linePrefix + "\u2502  " // "|  "
			}

			walkDir(sb_ptr, filepath.Join(root, file.Name()), nextLinePrefix)
		}
	}
}

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

func Run() {
	if err := _APP.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
