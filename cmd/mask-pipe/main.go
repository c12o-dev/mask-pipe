package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/c12o-dev/mask-pipe/internal/filter"
	"github.com/c12o-dev/mask-pipe/patterns"
)

var (
	version = "dev"
	commit  = "none"
)

const usage = `mask-pipe — filter secrets from terminal output

Usage:
  <command> | mask-pipe [flags]

Flags:
  -h, --help      show this help and exit
  -V, --version   show version and exit

Example:
  aws sts get-caller-identity 2>&1 | mask-pipe
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("mask-pipe", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var (
		showHelp    bool
		showVersion bool
	)
	fs.BoolVar(&showHelp, "help", false, "")
	fs.BoolVar(&showHelp, "h", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&showVersion, "V", false, "")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		fmt.Fprint(stderr, usage)
		return 2
	}

	switch {
	case showHelp:
		fmt.Fprint(stdout, usage)
		return 0
	case showVersion:
		fmt.Fprintf(stdout, "mask-pipe %s (%s)\n", version, commit)
		return 0
	}

	f := filter.New(patterns.Builtins, patterns.DefaultShowTail)
	if err := f.Run(stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		return 1
	}
	return 0
}
