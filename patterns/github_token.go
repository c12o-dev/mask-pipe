package patterns

import "regexp"

var githubToken = &Pattern{
	ID:         "github_token",
	Name:       "GitHub Token",
	Regex:      regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{36,}`),
	CaptureIdx: 0,
	Examples: []string{
		"ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"gho_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"ghu_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"ghs_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"ghr_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
	},
	NonExamples: []string{
		"ghp_short",
		"ghx_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"ghp-ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"GHP_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
		"ghprb_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh12",
		"prefixghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234",
	},
	Source: "https://github.blog/changelog/2021-03-31-authentication-token-format-updates-are-generally-available/",
}
