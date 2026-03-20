package laravel

import (
	"os"
	"path/filepath"
	"testing"
)

func makeDir(t *testing.T) string {
	t.Helper()
	return t.TempDir()
}

// ── IsLaravel ─────────────────────────────────────────────────────────────────

func TestIsLaravel_artisanFile(t *testing.T) {
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "artisan"), []byte("#!/usr/bin/env php\n"), 0755)
	if !IsLaravel(dir) {
		t.Error("expected true when artisan file is present")
	}
}

func TestIsLaravel_composerRequire(t *testing.T) {
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(`{
		"require": {
			"laravel/framework": "^11.0",
			"php": "^8.2"
		}
	}`), 0644)
	if !IsLaravel(dir) {
		t.Error("expected true when composer.json has laravel/framework in require")
	}
}

func TestIsLaravel_composerRequireDev(t *testing.T) {
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(`{
		"require": {"php": "^8.2"},
		"require-dev": {"laravel/framework": "^11.0"}
	}`), 0644)
	if !IsLaravel(dir) {
		t.Error("expected true when composer.json has laravel/framework in require-dev")
	}
}

func TestIsLaravel_publicIndexIlluminate(t *testing.T) {
	dir := makeDir(t)
	os.MkdirAll(filepath.Join(dir, "public"), 0755)
	os.WriteFile(filepath.Join(dir, "public", "index.php"), []byte(`<?php
require __DIR__.'/../vendor/autoload.php';
$app = require_once __DIR__.'/../bootstrap/app.php';
use Illuminate\Http\Request;
`), 0644)
	if !IsLaravel(dir) {
		t.Error("expected true when public/index.php contains Illuminate reference")
	}
}

func TestIsLaravel_notLaravel_emptyDir(t *testing.T) {
	dir := makeDir(t)
	if IsLaravel(dir) {
		t.Error("expected false for empty directory")
	}
}

func TestIsLaravel_notLaravel_composerNoFramework(t *testing.T) {
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(`{
		"require": {"guzzlehttp/guzzle": "^7.0"}
	}`), 0644)
	if IsLaravel(dir) {
		t.Error("expected false when composer.json has no laravel/framework")
	}
}

func TestIsLaravel_notLaravel_invalidComposerJSON(t *testing.T) {
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte(`not valid json`), 0644)
	if IsLaravel(dir) {
		t.Error("expected false for invalid composer.json")
	}
}

func TestIsLaravel_notLaravel_publicIndexNoIlluminate(t *testing.T) {
	dir := makeDir(t)
	os.MkdirAll(filepath.Join(dir, "public"), 0755)
	os.WriteFile(filepath.Join(dir, "public", "index.php"), []byte(`<?php echo "hello"; ?>`), 0644)
	if IsLaravel(dir) {
		t.Error("expected false when public/index.php has no Illuminate reference")
	}
}

func TestIsLaravel_artisanTakesPriority(t *testing.T) {
	// artisan present but also a bad composer.json — artisan wins
	dir := makeDir(t)
	os.WriteFile(filepath.Join(dir, "artisan"), []byte("#!/usr/bin/env php"), 0755)
	os.WriteFile(filepath.Join(dir, "composer.json"), []byte("not json"), 0644)
	if !IsLaravel(dir) {
		t.Error("expected true when artisan file is present, regardless of composer.json")
	}
}
