package patterns

import "regexp"

var awsAccessKey = &Pattern{
	ID:         "aws_access_key",
	Name:       "AWS Access Key ID",
	Regex:      regexp.MustCompile(`AKIA[0-9A-Z]{16}`),
	CaptureIdx: 0,
	Examples: []string{
		"AKIAIOSFODNN7EXAMPLE",
		"AKIAI44QH8DHBEXAMPLE",
		"AKIAJQABCDEFGHIJKLMN",
		"AKIAZZZZ9999AAAABBBB",
		"AKIA2RAY6KGQJM7Q3EYX",
	},
	NonExamples: []string{
		"AKIA",
		"AKIAIOSFODNN7EXAM",
		"XKIAIOSFODNN7EXAMPLE",
		"akiaiosfodnn7example",
		"AKIAIOSFODNN!EXAMPLE",
		"AKIAiosfodnn7example",
		"BKIAIOSFODNN7EXAMPLE",
	},
	Source: "https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_identifiers.html",
}
