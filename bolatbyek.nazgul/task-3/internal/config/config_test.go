package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Nazkaaa/task-3/internal/config"
)

const (
	testFilePerm = 0o600
	testDirPerm  = 0o755
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	configContent := `input-file: "data/input.xml"
output-file: "output/currencies.json"
`

	err := os.WriteFile(configPath, []byte(configContent), testFilePerm)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.InputFile != "data/input.xml" {
		t.Errorf("Expected InputFile to be 'data/input.xml', got '%s'", cfg.InputFile)
	}

	if cfg.OutputFile != "output/currencies.json" {
		t.Errorf("Expected OutputFile to be 'output/currencies.json', got '%s'", cfg.OutputFile)
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	t.Parallel()

	_, err := config.LoadConfig("non_existent_file.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid_config.yaml")
	configContent := `invalid: yaml: content`

	err := os.WriteFile(configPath, []byte(configContent), testFilePerm)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err = config.LoadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestEnsureOutputDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "subdir", "output.json")

	err := config.EnsureOutputDir(outputFile)
	if err != nil {
		t.Fatalf("EnsureOutputDir failed: %v", err)
	}

	// Check if directory was created
	dir := filepath.Dir(outputFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("Expected directory '%s' to exist, but it doesn't", dir)
	}
}

func TestEnsureOutputDir_ExistingDir(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.json")

	// Create directory first
	err := os.MkdirAll(tmpDir, testDirPerm)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	err = config.EnsureOutputDir(outputFile)
	if err != nil {
		t.Fatalf("EnsureOutputDir failed for existing directory: %v", err)
	}
}
