package main

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	// Create test XML file
	xmlPath := filepath.Join(tmpDir, "test_input.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<ValCurs Date="02.03.2022" name="Foreign Currency Market">
  <Valute ID="R01035">
    <NumCode>826</NumCode>
    <CharCode>GBP</CharCode>
    <Nominal>1</Nominal>
    <Name>Фунт стерлингов Соединенного королевства</Name>
    <Value>95,0000</Value>
  </Valute>
  <Valute ID="R01235">
    <NumCode>840</NumCode>
    <CharCode>USD</CharCode>
    <Nominal>1</Nominal>
    <Name>Доллар США</Name>
    <Value>120,0000</Value>
  </Valute>
</ValCurs>
`
	err := os.WriteFile(xmlPath, []byte(xmlContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	// Create test config file with absolute paths
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	outputPath := filepath.Join(tmpDir, "test_output.json")
	// Use forward slashes for YAML compatibility
	xmlPathNormalized := filepath.ToSlash(xmlPath)
	outputPathNormalized := filepath.ToSlash(outputPath)
	configContent := `input-file: "` + xmlPathNormalized + `"
output-file: "` + outputPathNormalized + `"
`
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Suppress stdout during test
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	err = run(configPath)

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	io.Copy(io.Discard, r)
	r.Close()

	if err != nil {
		t.Fatalf("run() failed: %v", err)
	}

	// Verify output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %v", err)
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	err := run("non_existent_config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent config file, got nil")
	}
}

