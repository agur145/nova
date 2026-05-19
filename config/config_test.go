package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDefaultsNovaDirToHomeNova(t *testing.T) {
	t.Setenv("NOVA_DIR", "")

	cfg := Load()
	want := normalizePath("~/.nova")
	if cfg.NovaDir != want {
		t.Fatalf("默认 NovaDir 不符合预期: want=%s got=%s", want, cfg.NovaDir)
	}
}

func TestLoadNovaDirFromEnv(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nova-data")
	t.Setenv("NOVA_DIR", dir)

	cfg := Load()
	if cfg.NovaDir != dir {
		t.Fatalf("环境变量 NovaDir 不符合预期: want=%s got=%s", dir, cfg.NovaDir)
	}
}

func TestNormalizePathExpandsRelativeAndHome(t *testing.T) {
	relative := "data/nova"
	abs, err := filepath.Abs(relative)
	if err != nil {
		t.Fatal(err)
	}
	if got := normalizePath(relative); got != abs {
		t.Fatalf("相对路径未转绝对路径: want=%s got=%s", abs, got)
	}

	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		t.Skip("当前环境无 home 目录")
	}
	want := filepath.Join(home, ".nova")
	if got := normalizePath("~/.nova"); got != want {
		t.Fatalf("~ 路径未正确展开: want=%s got=%s", want, got)
	}
}
