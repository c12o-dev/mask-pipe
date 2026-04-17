package patterns

import "regexp"

var jwt = &Pattern{
	ID:         "jwt",
	Name:       "JSON Web Token",
	Regex:      regexp.MustCompile(`eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+`),
	CaptureIdx: 0,
	Examples: []string{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJodHRwczovL2V4YW1wbGUuY29tIn0.signature123abc",
		"eyJhbGciOiJFUzI1NiJ9.eyJleHAiOjE2OTk5OTk5OTl9.MEUCIQD_abc123",
		"eyJhbGciOiJIUzM4NCJ9.eyJyb2xlIjoiYWRtaW4ifQ.abcdefghijklmnopqrstuvwx",
		"eyJhbGciOiJQUzUxMiIsImtpZCI6ImtleTEifQ.eyJhdWQiOiJhcGkuZXhhbXBsZS5jb20ifQ.sig_value_here",
	},
	NonExamples: []string{
		"eyJ",
		"eyJhbGciOiJIUzI1NiJ9",
		"eyJhbGciOiJIUzI1NiJ9.notbase64header",
		"eyXhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.sig",
		"eyJhbGciOiJIUzI1NiJ9.notEyJheader.signature",
	},
	Source: "https://datatracker.ietf.org/doc/html/rfc7519",
}
