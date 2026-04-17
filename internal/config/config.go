package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"regexp"
)

type Config struct {
	Builtin   map[string]bool  `toml:"builtin"`
	Custom    []CustomPattern  `toml:"custom"`
	Display   DisplayConfig    `toml:"display"`
	Allowlist []AllowlistEntry `toml:"allowlist"`
}

type CustomPattern struct {
	Name        string `toml:"name"`
	Pattern     string `toml:"pattern"`
	ShowTail    *int   `toml:"show_tail"`
	Replacement string `toml:"replacement"`
}

type DisplayConfig struct {
	MaskChar string `toml:"mask_char"`
	ShowTail int    `toml:"show_tail"`
	Color    bool   `toml:"color"`
}

type AllowlistEntry struct {
	Name    string `toml:"name"`
	Pattern string `toml:"pattern"`
}

func Default() *Config {
	return &Config{
		Builtin: nil,
		Display: DisplayConfig{
			MaskChar: "*",
			ShowTail: 4,
			Color:    true,
		},
	}
}

// Load finds and parses the config file. Returns Default() if no file found.
func Load(explicitPath string) (*Config, error) {
	path := findConfigFile(explicitPath)
	if path == "" {
		return Default(), nil
	}
	return loadFile(path)
}

func findConfigFile(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if env := os.Getenv("MASK_PIPE_CONFIG"); env != "" {
		return env
	}
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		p := filepath.Join(xdg, "mask-pipe", "config.toml")
		if fileExists(p) {
			return p
		}
	} else if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, ".config", "mask-pipe", "config.toml")
		if fileExists(p) {
			return p
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, ".mask-pipe.toml")
		if fileExists(p) {
			return p
		}
	}
	return ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func loadFile(path string) (*Config, error) {
	cfg := Default()
	if _, err := toml.DecodeFile(path, cfg); err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("config %s: %w", path, err)
	}
	return cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Display.MaskChar != "" {
		runes := []rune(cfg.Display.MaskChar)
		if len(runes) != 1 {
			return fmt.Errorf("[display].mask_char must be a single character, got %q", cfg.Display.MaskChar)
		}
	}

	names := make(map[string]bool)
	for i, c := range cfg.Custom {
		if c.Name == "" {
			return fmt.Errorf("[[custom]][%d]: name is required", i)
		}
		if c.Pattern == "" {
			return fmt.Errorf("[[custom]] %q: pattern is required", c.Name)
		}
		if _, err := regexp.Compile(c.Pattern); err != nil {
			return fmt.Errorf("[[custom]] %q: invalid regex: %w", c.Name, err)
		}
		if names[c.Name] {
			return fmt.Errorf("[[custom]] %q: duplicate name", c.Name)
		}
		names[c.Name] = true
	}

	for i, a := range cfg.Allowlist {
		if a.Name == "" {
			return fmt.Errorf("[[allowlist]][%d]: name is required", i)
		}
		if a.Pattern == "" {
			return fmt.Errorf("[[allowlist]] %q: pattern is required", a.Name)
		}
		if _, err := regexp.Compile(a.Pattern); err != nil {
			return fmt.Errorf("[[allowlist]] %q: invalid regex: %w", a.Name, err)
		}
	}
	return nil
}

// IsBuiltinEnabled checks if a built-in pattern is enabled. Default: true.
func (c *Config) IsBuiltinEnabled(id string) bool {
	if c.Builtin == nil {
		return true
	}
	enabled, exists := c.Builtin[id]
	if !exists {
		return true
	}
	return enabled
}

// MaskCharRune returns the mask character as a rune. Defaults to '*'.
func (c *Config) MaskCharRune() rune {
	if c.Display.MaskChar == "" {
		return '*'
	}
	return []rune(c.Display.MaskChar)[0]
}

// AllowlistRegexps compiles and returns allowlist patterns.
func (c *Config) AllowlistRegexps() []*regexp.Regexp {
	var res []*regexp.Regexp
	for _, a := range c.Allowlist {
		re, _ := regexp.Compile(a.Pattern) // already validated
		res = append(res, re)
	}
	return res
}

// CustomPatternRegexps returns compiled custom patterns ready for use.
func (c *Config) CustomPatterns() []CompiledCustom {
	var res []CompiledCustom
	for _, cp := range c.Custom {
		re, _ := regexp.Compile(cp.Pattern) // already validated
		cc := CompiledCustom{
			Name:        cp.Name,
			Regex:       re,
			Replacement: cp.Replacement,
			ShowTail:    c.Display.ShowTail,
		}
		if cp.ShowTail != nil {
			cc.ShowTail = *cp.ShowTail
		}
		res = append(res, cc)
	}
	return res
}

type CompiledCustom struct {
	Name        string
	Regex       *regexp.Regexp
	Replacement string
	ShowTail    int
}

// MaskChar returns the mask character string (for use with strings.Repeat).
func (c *Config) MaskCharStr() string {
	if c.Display.MaskChar == "" {
		return "*"
	}
	return string([]rune(c.Display.MaskChar)[:1])
}
