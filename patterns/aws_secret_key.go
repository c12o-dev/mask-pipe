package patterns

import "regexp"

var awsSecretKey = &Pattern{
	ID:         "aws_secret_key",
	Name:       "AWS Secret Access Key",
	Regex:      regexp.MustCompile(`(?i)aws_secret_access_key\s*[=:]\s*([A-Za-z0-9/+=]{40})`),
	CaptureIdx: 1,
	Examples: []string{
		"aws_secret_access_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"aws_secret_access_key:wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"AWS_SECRET_ACCESS_KEY : Ab0dEfGhIjKlMnOpQrStUvWxYz012345678+/A==",
	},
	NonExamples: []string{
		"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"aws_secret_access_key=short",
		"aws_access_key_id=AKIAIOSFODNN7EXAMPLE",
		"some_other_key=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		"aws_secret_access_key",
	},
	Source: "https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html",
}
