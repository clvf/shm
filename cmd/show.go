package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v2"
)

func show(cCtx *cli.Context) error {
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
			glamour.WithColorProfile(termenv.TrueColor),
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
			fmt.Fprintf(sb_ptr, "%v%v%v%v%v\n", linePrefix, filePrefix, _COLOR_BLUE, file.Name(), _COLOR_RESET)
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
