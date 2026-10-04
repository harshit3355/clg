package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadTimezone(t *testing.T) {
	t.Run("defaults to UTC", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "")
		cfg, err := Load(t.TempDir(), t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Timezone != "UTC" {
			t.Fatalf("Timezone = %q, want UTC", cfg.Timezone)
		}
	})

	t.Run("project config overrides user config", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "")
		homeDir := t.TempDir()
		projectDir := t.TempDir()
		writeConfig(t, homeDir, "timezone: Europe/Paris\n")
		writeConfig(t, projectDir, "timezone: Asia/Kolkata\n")

		cfg, err := Load(homeDir, projectDir)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Timezone != "Asia/Kolkata" {
			t.Fatalf("Timezone = %q, want Asia/Kolkata", cfg.Timezone)
		}
	})

	t.Run("environment overrides project config", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "America/Los_Angeles")
		projectDir := t.TempDir()
		writeConfig(t, projectDir, "timezone: Asia/Kolkata\n")

		cfg, err := Load(t.TempDir(), projectDir)
		if err != nil {
			t.Fatal(err)
		}
		if cfg.Timezone != "America/Los_Angeles" {
			t.Fatalf("Timezone = %q, want America/Los_Angeles", cfg.Timezone)
		}
	})

	t.Run("rejects an unknown timezone", func(t *testing.T) {
		t.Setenv("CLG_TIMEZONE", "")
		projectDir := t.TempDir()
		writeConfig(t, projectDir, "timezone: Mars/Olympus\n")

		_, err := Load(t.TempDir(), projectDir)
		if err == nil || !strings.HasPrefix(err.Error(), `invalid timezone "Mars/Olympus":`) {
			t.Fatalf("Load() error = %v, want invalid timezone error", err)
		}
	})
}

func writeConfig(t *testing.T, dir, contents string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, configFilename), []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
