package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.LargeFiles.MinSizeMB != 100 {
		t.Errorf("expected MinSizeMB=100, got %d", cfg.LargeFiles.MinSizeMB)
	}
}

func TestLoadConfig_FromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.toml")

	content := `
[large_files]
min_size_mb = 200
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	os.Setenv("SLIMMING_CONFIG", configPath)
	defer os.Unsetenv("SLIMMING_CONFIG")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.LargeFiles.MinSizeMB != 200 {
		t.Errorf("expected MinSizeMB=200, got %d", cfg.LargeFiles.MinSizeMB)
	}
}
