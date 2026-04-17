package patterns

import "regexp"

var pemPrivateKey = &Pattern{
	ID:          "pem_private_key",
	Name:        "PEM private key block",
	Regex:       regexp.MustCompile(`(?s)-----BEGIN [A-Z ]*PRIVATE KEY-----.*?-----END [A-Z ]*PRIVATE KEY-----`),
	CaptureIdx:  0,
	Replacement: "[REDACTED PRIVATE KEY]",
	Examples: []string{
		"-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJBALRiMLAH\n-----END RSA PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----\nMHQCAQEEIBkg\n-----END EC PRIVATE KEY-----",
		"-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIA\n-----END PRIVATE KEY-----",
		"-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjE\n-----END OPENSSH PRIVATE KEY-----",
		"-----BEGIN ENCRYPTED PRIVATE KEY-----\nMIIFDjBABgkqhkiG9w0B\n-----END ENCRYPTED PRIVATE KEY-----",
	},
	NonExamples: []string{
		"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0B\n-----END PUBLIC KEY-----",
		"-----BEGIN CERTIFICATE-----\nMIICpDCCAYwCCQC7RXjA\n-----END CERTIFICATE-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----END RSA PRIVATE KEY-----",
		"The key format uses -----BEGIN PRIVATE KEY----- markers.",
	},
	Source:      "https://datatracker.ietf.org/doc/html/rfc7468",
	Multiline:   true,
	BeginMarker: regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----`),
	EndMarker:   regexp.MustCompile(`-----END [A-Z ]*PRIVATE KEY-----`),
}
