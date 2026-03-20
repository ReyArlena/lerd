package certs

import (
	"os"
	"path/filepath"
	"testing"
)

// ── MkcertPath ────────────────────────────────────────────────────────────────

func TestMkcertPath_usesDataDir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)
	got := MkcertPath()
	want := filepath.Join(tmp, "lerd", "bin", "mkcert")
	if got != want {
		t.Errorf("MkcertPath() = %q, want %q", got, want)
	}
}

// ── CertExists ────────────────────────────────────────────────────────────────

func TestCertExists_returnsFalseWhenMissing(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)
	if CertExists("myapp.test") {
		t.Error("expected false for non-existent cert")
	}
}

func TestCertExists_returnsTrueWhenPresent(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	// Create the expected cert file path
	certsDir := filepath.Join(tmp, "lerd", "certs", "sites")
	os.MkdirAll(certsDir, 0755)
	os.WriteFile(filepath.Join(certsDir, "myapp.test.crt"), []byte("fake cert"), 0644)

	if !CertExists("myapp.test") {
		t.Error("expected true when cert file exists")
	}
}

func TestCertExists_onlyCrtRequired(t *testing.T) {
	// CertExists checks for .crt only, not .key
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	certsDir := filepath.Join(tmp, "lerd", "certs", "sites")
	os.MkdirAll(certsDir, 0755)
	// .crt exists, no .key
	os.WriteFile(filepath.Join(certsDir, "site.test.crt"), []byte("fake cert"), 0644)

	if !CertExists("site.test") {
		t.Error("expected true when only .crt file exists")
	}
}
