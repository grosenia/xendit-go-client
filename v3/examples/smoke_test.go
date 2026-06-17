package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Smoke test CLI examples — versi gagal (tanpa network).
// Harapan: exit code 1 + pesan jelas di stderr/stdout.

func examplesDir(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(dir) != "examples" {
		t.Skip("run from v3/examples directory")
	}
	return dir
}

func runExample(t *testing.T, dir, command string, extraEnv ...string) (string, int) {
	t.Helper()
	cmd := exec.Command("go", "run", ".", command)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), extraEnv...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return string(out), 0
	}
	if exit, ok := err.(*exec.ExitError); ok {
		return string(out), exit.ExitCode()
	}
	t.Fatalf("run %s: %v", command, err)
	return string(out), 1
}

func writeTempConfig(t *testing.T, dir string, lines ...string) {
	t.Helper()
	path := filepath.Join(dir, "smoke-test.props")
	content := strings.Join(lines, "\n") + "\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(path) })
}

func TestCLI_UnknownCommand_Fail(t *testing.T) {
	dir := examplesDir(t)
	out, code := runExample(t, dir, "not-a-command")
	if code == 0 {
		t.Fatalf("expected exit 1, got 0\n%s", out)
	}
	if !strings.Contains(out, "unknown command") {
		t.Fatalf("expected unknown command message\n%s", out)
	}
}

func TestCLI_PaymentRequest_MissingToken_Fail(t *testing.T) {
	dir := examplesDir(t)
	writeTempConfig(t, dir,
		"KEY_WRITE_MONEY_IN=xnd_development_fake",
		"API_VERSION=2024-11-11",
		"PAYMENT_TOKEN_ID=",
		"CVN=",
		"AMOUNT=50000",
	)
	out, code := runExample(t, dir, "payment-request", "XENDIT_EXAMPLE_CONFIG=smoke-test")
	if code == 0 {
		t.Fatalf("expected exit 1, got 0\n%s", out)
	}
	if !strings.Contains(out, "PAYMENT_TOKEN_ID") {
		t.Fatalf("expected PAYMENT_TOKEN_ID error\n%s", out)
	}
	if !strings.Contains(out, "GAGAL") {
		t.Fatalf("expected HASIL GAGAL summary\n%s", out)
	}
}

func TestCLI_GetToken_MissingToken_Fail(t *testing.T) {
	dir := examplesDir(t)
	writeTempConfig(t, dir,
		"KEY_WRITE_MONEY_IN=xnd_development_fake",
		"API_VERSION=2024-11-11",
		"PAYMENT_TOKEN_ID=REPLACE_ME",
	)
	out, code := runExample(t, dir, "get-token", "XENDIT_EXAMPLE_CONFIG=smoke-test")
	if code == 0 {
		t.Fatalf("expected exit 1, got 0\n%s", out)
	}
	if !strings.Contains(out, "PAYMENT_TOKEN_ID") {
		t.Fatalf("expected PAYMENT_TOKEN_ID error\n%s", out)
	}
}

func TestCLI_PrintUsage_NoArgs_Fail(t *testing.T) {
	dir := examplesDir(t)
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected exit 1\n%s", string(out))
	}
	if !strings.Contains(string(out), "card-session") {
		t.Fatalf("expected usage with commands\n%s", string(out))
	}
}

func TestLogHelpers_MaskKey(t *testing.T) {
	if got := maskKey(""); got != "(empty/invalid)" {
		t.Fatalf("empty key: %q", got)
	}
	if got := maskKey("xnd_development_ABCDEFGHIJK"); !strings.Contains(got, "...") {
		t.Fatalf("expected masked key: %q", got)
	}
}

func TestLogHelpers_ApiErr(t *testing.T) {
	err := apiErr("[CODE] message")
	if err == nil || err.Error() != "[CODE] message" {
		t.Fatalf("unexpected: %v", err)
	}
}
