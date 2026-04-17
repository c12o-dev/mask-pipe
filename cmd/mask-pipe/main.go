package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/tabwriter"

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
  mask-pipe <subcommand>

Flags:
  -h, --help            show this help and exit
  -V, --version         show version and exit
      --config <path>   path to config file (default: auto-detect)
      --mask-char <c>   override masking character
      --show-tail <N>   show last N chars of masked values (0 = full mask)

Subcommands:
  list-patterns   print all active built-in and custom patterns
  doctor          diagnose configuration and patterns
  version         print version and build info

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

	// Subcommand routing
	if fs.NArg() > 0 {
		switch fs.Arg(0) {
		case "list-patterns":
			return cmdListPatterns(configPath, stdout, stderr)
		case "doctor":
			return cmdDoctor(configPath, stdout, stderr)
		case "version":
			fmt.Fprintf(stdout, "mask-pipe %s (%s)\n", version, commit)
			return 0
		default:
			fmt.Fprintf(stderr, "mask-pipe: unknown subcommand %q\n", fs.Arg(0))
			fmt.Fprint(stderr, usage)
			return 2
		}
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		return 1
	}

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

func cmdListPatterns(configPath string, stdout, stderr io.Writer) int {
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stderr, "mask-pipe: %v\n", err)
		return 1
	}

	w := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSOURCE\tSTATUS\tREGEX")

	for _, p := range patterns.Builtins {
		status := "enabled"
		if !cfg.IsBuiltinEnabled(p.ID) {
			status = "disabled"
		}
		fmt.Fprintf(w, "%s\tbuiltin\t%s\t%s\n", p.ID, status, p.Regex.String())
	}

	for _, cp := range cfg.CustomPatterns() {
		fmt.Fprintf(w, "custom:%s\tcustom\tenabled\t%s\n", cp.Name, cp.Regex.String())
	}

	w.Flush()
	return 0
}

func cmdDoctor(configPath string, stdout, stderr io.Writer) int {
	ok := true

	// Check config
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(stdout, "x Config: %v\n", err)
		ok = false
	} else if configPath != "" {
		fmt.Fprintf(stdout, "ok Config: loaded %s\n", configPath)
	} else {
		fmt.Fprintln(stdout, "ok Config: using defaults (no config file found)")
	}

	if cfg == nil {
		fmt.Fprintln(stdout, "\nDoctor found problems.")
		return 1
	}

	// Check built-in patterns
	enabled := 0
	for _, p := range patterns.Builtins {
		if cfg.IsBuiltinEnabled(p.ID) {
			enabled++
		}
	}
	fmt.Fprintf(stdout, "ok Built-in patterns: %d/%d enabled\n", enabled, len(patterns.Builtins))

	// Check custom patterns
	customs := cfg.CustomPatterns()
	if len(customs) > 0 {
		fmt.Fprintf(stdout, "ok Custom patterns: %d defined, all regexes compile\n", len(customs))
	} else {
		fmt.Fprintln(stdout, "ok Custom patterns: 0 defined")
	}

	// Check allowlist
	allowlist := cfg.AllowlistRegexps()
	if len(allowlist) > 0 {
		fmt.Fprintf(stdout, "ok Allowlist: %d entries\n", len(allowlist))
	}

	// Check stdout writable
	if _, err := fmt.Fprint(stdout, ""); err != nil {
		fmt.Fprintln(stdout, "x stdout: not writable")
		ok = false
	} else {
		fmt.Fprintln(stdout, "ok stdout: writable")
	}

	if ok {
		fmt.Fprintln(stdout, "\nAll checks passed.")
		return 0
	}
	fmt.Fprintln(stdout, "\nDoctor found problems.")
	return 1
}

func buildPatterns(cfg *config.Config) []*patterns.Pattern {
	var pats []*patterns.Pattern

	for _, p := range patterns.Builtins {
		if cfg.IsBuiltinEnabled(p.ID) {
			pats = append(pats, p)
		}
	}

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
