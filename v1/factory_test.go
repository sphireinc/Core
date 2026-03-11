package core

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFindConfigFilePrefersExplicitEnvVar(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "custom.json")
	if err := os.WriteFile(configPath, []byte(`{}`), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Setenv(ConfigPathEnvironmentVariableName, configPath)

	got, err := findConfigFile(configLocations(), "dev.json")
	if err != nil {
		t.Fatal(err)
	}

	if got != configPath {
		t.Fatalf("expected %q, got %q", configPath, got)
	}
}

func TestLoadSupportsLegacyAndNormalizedKeys(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "dev.json")

	raw := `{
		"application": {"name": "core"},
		"components": {
			"bigCache": false,
			"memCache": false,
			"redis": false,
			"mySQL": false,
			"HTTPCache": true,
			"log": false
		},
		"server": {
			"address": "127.0.0.1",
			"port": "8080",
			"writeTimeout": 10,
			"read_timeout": "5s",
			"memCacheTime": "30s"
		},
		"middleware": {
			"limiter": {"enabled": true, "time": "3s", "max": 9}
		},
		"persistence": {
			"in_memory": {"capacity": 321, "ttl": "45s"},
			"neo4J": {}
		}
	}`

	if err := os.WriteFile(configPath, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg := &Config{
		Environment: Environment{Location: configPath},
	}
	if err := cfg.Load(); err != nil {
		t.Fatal(err)
	}

	if !cfg.Components.HTTPCache {
		t.Fatal("expected http_cache to be true after key normalization")
	}
	if got := cfg.Server.WriteTimeout.Or(0); got != 10*time.Second {
		t.Fatalf("expected write timeout 10s, got %s", got)
	}
	if got := cfg.Server.ReadTimeout.Or(0); got != 5*time.Second {
		t.Fatalf("expected read timeout 5s, got %s", got)
	}
	if got := cfg.Server.MemCacheTTL.Or(0); got != 30*time.Second {
		t.Fatalf("expected mem cache ttl 30s, got %s", got)
	}
	if got := cfg.Persistence.InMemory.Capacity; got != 321 {
		t.Fatalf("expected in_memory capacity 321, got %d", got)
	}
	if got := cfg.Persistence.InMemory.TTL.Or(0); got != 45*time.Second {
		t.Fatalf("expected in_memory ttl 45s, got %s", got)
	}
}
