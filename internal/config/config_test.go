package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.json")
	data, _ := json.Marshal(Config{
		Lark: LarkConfig{AppID: "file_id", AppSecret: "file_secret"},
	})
	os.WriteFile(cfgPath, data, 0600)
	t.Setenv("HOME", dir)
	os.MkdirAll(filepath.Join(dir, ".config", "chathub"), 0700)
	os.WriteFile(filepath.Join(dir, ".config", "chathub", "config.json"), data, 0600)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lark.AppID != "file_id" {
		t.Errorf("got %q, want file_id", cfg.Lark.AppID)
	}
}

func TestEnvOverridesFile(t *testing.T) {
	dir := t.TempDir()
	data, _ := json.Marshal(Config{
		Lark: LarkConfig{AppID: "file_id", AppSecret: "file_secret"},
	})
	os.MkdirAll(filepath.Join(dir, ".config", "chathub"), 0700)
	os.WriteFile(filepath.Join(dir, ".config", "chathub", "config.json"), data, 0600)
	t.Setenv("HOME", dir)
	t.Setenv("CHATHUB_LARK_APP_ID", "env_id")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Lark.AppID != "env_id" {
		t.Errorf("got %q, want env_id", cfg.Lark.AppID)
	}
	if cfg.Lark.AppSecret != "file_secret" {
		t.Errorf("got %q, want file_secret", cfg.Lark.AppSecret)
	}
}
