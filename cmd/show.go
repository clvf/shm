package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/glamour"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/muesli/termenv"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
)

func show(cCtx *cli.Context) error {
	if !cCtx.Args().Present() {
		printRepo(_REPO)
		return nil
	}

	file, err := filepath.Abs(filepath.Join(_REPO, cCtx.Args().First()))
	check(err)

	info, err := os.Stat(file)
	check(err)
	if info.IsDir() {
		printRepo(file)
		return nil
	}

	content, err := os.ReadFile(file)
	check(err)

	stdoutFd := int(os.Stdout.Fd())

	if !term.IsTerminal(stdoutFd) {
		os.Stdout.Write(content)
		return nil
	}

	termWidth, termHeight, err := term.GetSize(stdoutFd)
	check(err)

	out := string(content)

	if !_MONOCHROME {
		style := _DEFAULT_STYLE
		if override := os.Getenv("GLAMOUR_STYLE"); override != "" {
			style = override
		}

		r, err := glamour.NewTermRenderer(
			glamour.WithStandardStyle(style),
			glamour.WithWordWrap(termWidth),
			glamour.WithColorProfile(termenv.TrueColor),
		)
		check(err)

		out, err = r.Render(string(content))
		check(err)
	}

	lineCnt := strings.Count(out, "\n")

	if termHeight <= lineCnt {
		pager := os.Getenv("PAGER")
		if pager == "" {
			pager = "/bin/less -iR"
		}

		penv, args, err := shellwords.ParseWithEnvs(pager)

		binary, err := exec.LookPath(args[0])
		check(err)

		f, err := os.CreateTemp("", "shm.*.txt")
		check(err)
		// NOTE: no defer (remove) as we're exec'ing into "pager" so it
		//       wouldn't be called anyway.

		_, err = f.Write([]byte(out))
		check(err)

		err = f.Close()
		check(err)

		env := append(os.Environ(), penv...)
		args = append(args, f.Name())

		syscall.Exec(binary, args, env)
	}

	fmt.Println(out)
	return nil
}

func printRepo(basedir string) {
	fmt.Printf("%v%v%v\n", _COLOR_BLUE, filepath.Base(basedir), _COLOR_RESET)

	var sb strings.Builder
	walkDir(&sb, basedir, "")
	fmt.Print(sb.String())
}

func walkDir(sb_ptr *strings.Builder, root string, linePrefix string) {
	files, err := os.ReadDir(root)
	check(err)

	fileNum := len(files)
	for i, file := range files {
		filePrefix := "\u2514\u2500\u2500 " // "|-- "

		if i < fileNum-1 {
			filePrefix = "\u251c\u2500\u2500 " // "`-- "
		}

		// if it's just a file
		if !file.IsDir() {
			fmt.Fprintf(sb_ptr, "%v%v%v\n", linePrefix, filePrefix, file.Name())
			continue
		}

		// it's a directory then
		if _MONOCHROME {
			fmt.Fprintf(sb_ptr, "%v%v%v\n", linePrefix, filePrefix, file.Name())
		} else {
			fmt.Fprintf(sb_ptr, "%v%v%v%v%v\n", linePrefix, filePrefix, _COLOR_BLUE, file.Name(), _COLOR_RESET)
		}

		nextLinePrefix := linePrefix + "    "
		if i < fileNum-1 {
			nextLinePrefix = linePrefix + "\u2502   " // "|   "
		}

		walkDir(sb_ptr, filepath.Join(root, file.Name()), nextLinePrefix)
	}
}
