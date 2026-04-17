package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/c12o-dev/mask-pipe/internal/config"
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
  -h, --help            show this help and exit
  -V, --version         show version and exit
      --config <path>   path to config file (default: auto-detect)
      --mask-char <c>   override masking character
      --show-tail <N>   show last N chars of masked values (0 = full mask)

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
		configPath  string
		maskChar    string
		showTail    int
	)
	fs.BoolVar(&showHelp, "help", false, "")
	fs.BoolVar(&showHelp, "h", false, "")
	fs.BoolVar(&showVersion, "version", false, "")
	fs.BoolVar(&showVersion, "V", false, "")
	fs.StringVar(&configPath, "config", "", "")
	fs.StringVar(&maskChar, "mask-char", "", "")
	fs.IntVar(&showTail, "show-tail", -1, "")

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

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		return 1
	}

	// CLI flags override config
	if maskChar != "" {
		cfg.Display.MaskChar = maskChar
	}
	if showTail >= 0 {
		cfg.Display.ShowTail = showTail
	}

	pats := buildPatterns(cfg)

	f := &filter.Filter{
		Patterns:  pats,
		ShowTail:  cfg.Display.ShowTail,
		MaskChar:  cfg.MaskCharStr(),
		Allowlist: cfg.AllowlistRegexps(),
		Stderr:    stderr,
	}

	if err := f.Run(stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		return 1
	}
	return 0
}

func buildPatterns(cfg *config.Config) []*patterns.Pattern {
	var pats []*patterns.Pattern

	// Built-in patterns (filtered by config)
	for _, p := range patterns.Builtins {
		if cfg.IsBuiltinEnabled(p.ID) {
			pats = append(pats, p)
		}
	}

	// Custom patterns from config
	for _, cp := range cfg.CustomPatterns() {
		p := &patterns.Pattern{
			ID:          "custom:" + cp.Name,
			Name:        cp.Name,
			Regex:       cp.Regex,
			CaptureIdx:  0,
			Replacement: cp.Replacement,
		}
		pats = append(pats, p)
	}

	return pats
}
