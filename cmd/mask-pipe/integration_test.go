package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "mask-pipe-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binaryPath = filepath.Join(dir, "mask-pipe")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

func runBinary(t *testing.T, stdin string, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	cmd.Stdin = strings.NewReader(stdin)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			t.Fatalf("exec error: %v", err)
		}
	}
	return outBuf.String(), errBuf.String(), exitCode
}

func TestIntegrationMaskAWSKey(t *testing.T) {
	out, _, code := runBinary(t, "token AKIAIOSFODNN7EXAMPLE found\n")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "AKIA") || strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("expected masked AWS key, got: %q", out)
	}
}

func TestIntegrationMaskJWT(t *testing.T) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U"
	out, _, code := runBinary(t, "Bearer "+jwt+"\n")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if strings.Contains(out, jwt) {
		t.Errorf("JWT should be masked, got: %q", out)
	}
	if !strings.Contains(out, "Bearer ") {
		t.Errorf("Bearer prefix should be preserved, got: %q", out)
	}
}

func TestIntegrationMaskDBURL(t *testing.T) {
	out, _, code := runBinary(t, "postgres://admin:s3cretP4ss@db.example.com:5432/mydb\n")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if strings.Contains(out, "s3cretP4ss") {
		t.Errorf("password should be masked, got: %q", out)
	}
	if !strings.Contains(out, "****@") {
		t.Errorf("expected literal **** replacement, got: %q", out)
	}
}

func TestIntegrationMaskPEM(t *testing.T) {
	pem := "-----BEGIN RSA PRIVATE KEY-----\nMIIBogIBAAJBALRi\nbase64==\n-----END RSA PRIVATE KEY-----\n"
	out, _, code := runBinary(t, "before\n"+pem+"after\n")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "[REDACTED PRIVATE KEY]") {
		t.Errorf("PEM block should be redacted, got: %q", out)
	}
	if !strings.Contains(out, "before\n") || !strings.Contains(out, "after\n") {
		t.Errorf("surrounding text should be preserved, got: %q", out)
	}
}

func TestIntegrationCleanPassthrough(t *testing.T) {
	input := "no secrets here\njust normal output\n"
	out, _, code := runBinary(t, input)
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if out != input {
		t.Errorf("clean input should pass through unchanged: got %q, want %q", out, input)
	}
}

func TestIntegrationDryRun(t *testing.T) {
	out, _, code := runBinary(t, "AKIAIOSFODNN7EXAMPLE\n", "--dry-run", "--no-color")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "[MATCH:aws_access_key]") {
		t.Errorf("dry-run should show MATCH tags, got: %q", out)
	}
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("dry-run should preserve original value, got: %q", out)
	}
}

func TestIntegrationShowTailZero(t *testing.T) {
	out, _, code := runBinary(t, "AKIAIOSFODNN7EXAMPLE\n", "--show-tail", "0")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if strings.Contains(out, "MPLE") {
		t.Errorf("show-tail=0 should fully mask, got: %q", out)
	}
}

func TestIntegrationConfigBuiltinToggle(t *testing.T) {
	cfgDir := t.TempDir()
	cfgPath := filepath.Join(cfgDir, "config.toml")
	os.WriteFile(cfgPath, []byte("[builtin]\naws_access_key = false\n"), 0644)

	out, _, code := runBinary(t, "AKIAIOSFODNN7EXAMPLE\n", "--config", cfgPath)
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("disabled pattern should pass through, got: %q", out)
	}
}

func TestIntegrationVersionFlag(t *testing.T) {
	out, _, code := runBinary(t, "", "--version")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "mask-pipe") {
		t.Errorf("version output should contain 'mask-pipe', got: %q", out)
	}
}

func TestIntegrationDoctorSubcommand(t *testing.T) {
	out, _, code := runBinary(t, "", "doctor")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "All checks passed") {
		t.Errorf("doctor should pass, got: %q", out)
	}
}

func TestIntegrationListPatterns(t *testing.T) {
	out, _, code := runBinary(t, "", "list-patterns")
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if !strings.Contains(out, "aws_access_key") {
		t.Errorf("list-patterns should show aws_access_key, got: %q", out)
	}
	if strings.Count(out, "builtin") < 8 {
		t.Errorf("expected at least 8 builtin patterns, got: %q", out)
	}
}

func TestIntegrationInvalidFlag(t *testing.T) {
	_, _, code := runBinary(t, "", "--bogus")
	if code != 2 {
		t.Errorf("invalid flag: exit %d, want 2", code)
	}
}

func TestIntegrationMultipleSecrets(t *testing.T) {
	input := "AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE\nGITHUB_TOKEN=ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234\nclean line\n"
	out, _, code := runBinary(t, input)
	if code != 0 {
		t.Fatalf("exit %d", code)
	}
	if strings.Contains(out, "AKIAIOSFODNN7EXAMPLE") {
		t.Errorf("AWS key should be masked")
	}
	if strings.Contains(out, "ghp_ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh1234") {
		t.Errorf("GitHub token should be masked")
	}
	if !strings.Contains(out, "clean line") {
		t.Errorf("clean line should pass through")
	}
}
