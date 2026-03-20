package node

import (
	"os"
	"path/filepath"
	"testing"
)

// ── isNumericVersion ─────────────────────────────────────────────────────────

func TestIsNumericVersion(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"22", true},
		{"18", true},
		{"", false},
		{"system", false},
		{"lts/iron", false},
		{"v22", false},
		{"22.x", false},
		{"22.1.0", false},
	}
	for _, c := range cases {
		got := isNumericVersion(c.in)
		if got != c.want {
			t.Errorf("isNumericVersion(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// ── extractMajor ─────────────────────────────────────────────────────────────

func TestExtractMajor(t *testing.T) {
	cases := []struct{ in, want string }{
		{"22", "22"},
		{"18.12.0", "18"},
		{"20.1", "20"},
		{"system", "system"},
		{"", ""},
	}
	for _, c := range cases {
		got := extractMajor(c.in)
		if got != c.want {
			t.Errorf("extractMajor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ── parseNodeConstraint ──────────────────────────────────────────────────────

func TestParseNodeConstraint(t *testing.T) {
	cases := []struct{ in, want string }{
		{">=18", "18"},
		{"^20.0.0", "20"},
		{"18.x", "18"},
		{">=16.0.0 <20", "16"},
		{"*", ""},
		{"", ""},
	}
	for _, c := range cases {
		got := parseNodeConstraint(c.in)
		if got != c.want {
			t.Errorf("parseNodeConstraint(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// ── DetectVersion ────────────────────────────────────────────────────────────

func TestDetectVersion_nvmrc(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("v18.12.0\n"), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "18" {
		t.Errorf("got %q, want %q", got, "18")
	}
}

func TestDetectVersion_nvmrc_majorOnly(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("20\n"), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "20" {
		t.Errorf("got %q, want %q", got, "20")
	}
}

// Regression: .nvmrc containing "system" should fall through, not propagate "system".
func TestDetectVersion_nvmrc_system_fallsThrough(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("system\n"), 0644)
	// No .node-version, no package.json → should reach global default ("22")
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	// "system" is non-numeric and must not be returned
	if got == "system" {
		t.Error("DetectVersion must not return \"system\" from .nvmrc")
	}
}

func TestDetectVersion_nodeVersion_file(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".node-version"), []byte("v20.5.0\n"), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "20" {
		t.Errorf("got %q, want %q", got, "20")
	}
}

func TestDetectVersion_nodeVersion_precedence(t *testing.T) {
	// .nvmrc says "system" (invalid), .node-version says 16 → should use 16
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("system\n"), 0644)
	os.WriteFile(filepath.Join(dir, ".node-version"), []byte("16\n"), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "16" {
		t.Errorf("got %q, want %q", got, "16")
	}
}

func TestDetectVersion_packageJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"engines":{"node":">=18.0.0"}}`), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "18" {
		t.Errorf("got %q, want %q", got, "18")
	}
}

func TestDetectVersion_noFiles_returnsDefault(t *testing.T) {
	dir := t.TempDir()
	// XDG_CONFIG_HOME points to a dir with no config.yaml → defaultConfig() → "22"
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got != "22" {
		t.Errorf("got %q, want %q", got, "22")
	}
}

func TestDetectVersion_nvmrcOverridesPackageJSON(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, ".nvmrc"), []byte("20\n"), 0644)
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(`{"engines":{"node":">=18"}}`), 0644)
	got, err := DetectVersion(dir)
	if err != nil {
		t.Fatal(err)
	}
	// .nvmrc has priority over package.json
	if got != "20" {
		t.Errorf("got %q, want %q", got, "20")
	}
}
