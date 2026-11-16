package main

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nesoriel/neurolink/playerapi"
	"github.com/Nesoriel/neurolink/statusapi"
	"github.com/Nesoriel/neurolink/tui"
)

func TestParseAppConfigUsesEnvAPIKey(t *testing.T) {
	cfg, err := parseAppConfig(withConfigPath(t, nil), testLookupEnv(map[string]string{
		envAPIKey: "env-key",
	}))
	if err != nil {
		t.Fatalf("parseAppConfig() error = %v", err)
	}
	if cfg.APIKey != "env-key" {
		t.Fatalf("API key = %q, want env-key", cfg.APIKey)
	}
	if cfg.PollInterval != time.Minute {
		t.Fatalf("poll interval = %s, want 1m", cfg.PollInterval)
	}
	if cfg.Language != tui.LanguageEnglish {
		t.Fatalf("language = %s, want en", cfg.Language)
	}
}

func TestConfigSetKeyValueCommand(t *testing.T) {
	path := testConfigPath(t)
	var out bytes.Buffer
	if err := runConfigCommand([]string{"config", "set", "--config", path, "poll-interval", "40s"}, nil, &out); err != nil {
		t.Fatalf("runConfigCommand(set key value) error = %v", err)
	}

	loaded, err := loadFileConfig(path)
	if err != nil {
		t.Fatalf("loadFileConfig() error = %v", err)
	}
	if loaded.PollInterval != "40s" {
		t.Fatalf("poll interval = %q, want 40s", loaded.PollInterval)
	}
}

func TestParseAppConfigRejectsInvalidLanguage(t *testing.T) {
	_, err := parseAppConfig(withConfigPath(t, []string{"--lang", "fr"}), nil)
	if err == nil {
		t.Fatal("expected error for invalid language")
	}
}

func TestParseAppConfigRejectsInvalidInterval(t *testing.T) {
	_, err := parseAppConfig(withConfigPath(t, []string{"--poll-interval", "0s"}), nil)
	if err == nil {
		t.Fatal("expected error for zero poll interval")
	}
}

func TestParseAppConfigRejectsInvalidConfigInterval(t *testing.T) {
	path := testConfigPath(t)
	if err := saveFileConfig(path, fileConfig{PollInterval: "0s"}); err != nil {
		t.Fatalf("saveFileConfig() error = %v", err)
	}

	_, err := parseAppConfig([]string{"--config", path}, nil)
	if err == nil {
		t.Fatal("expected error for zero config poll interval")
	}
}

func TestAppConfigProviderSelectsDemoWithoutKey(t *testing.T) {
	cfg := appConfig{}
	snapshot, err := cfg.provider().Fetch(context.Background())
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if snapshot.Source != statusapi.SourceDemo {
		t.Fatalf("source = %s, want demo", snapshot.Source)
	}
}

func TestAppConfigPlayerProviderSelectsDemoWithoutKey(t *testing.T) {
	cfg := appConfig{}
	snapshot, err := cfg.playerProvider().Lookup(context.Background(), playerapi.LookupRequest{Player: "Demo", Platform: playerapi.PlatformPC})
	if err != nil {
		t.Fatalf("Lookup() error = %v", err)
	}
	if snapshot.Source != playerapi.SourceDemo {
		t.Fatalf("source = %s, want demo", snapshot.Source)
	}
}

func withConfigPath(t *testing.T, args []string) []string {
	t.Helper()
	result := []string{"--config", testConfigPath(t)}
	return append(result, args...)
}

func testConfigPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "neurolink", "config.json")
}

func testLookupEnv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
