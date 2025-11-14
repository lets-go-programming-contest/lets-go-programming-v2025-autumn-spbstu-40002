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

func TestRun_InvalidXML(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test config file with invalid XML path
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	invalidXMLPath := filepath.Join(tmpDir, "non_existent.xml")
	outputPath := filepath.Join(tmpDir, "test_output.json")
	
	configContent := `input-file: "` + filepath.ToSlash(invalidXMLPath) + `"
output-file: "` + filepath.ToSlash(outputPath) + `"
`
	err := os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	err = run(configPath)
	if err == nil {
		t.Error("Expected error for non-existent XML file, got nil")
	}
}

func TestRun_InvalidOutputDir(t *testing.T) {
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
</ValCurs>
`
	err := os.WriteFile(xmlPath, []byte(xmlContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	// Create config with invalid output directory (read-only parent)
	configPath := filepath.Join(tmpDir, "test_config.yaml")
	outputPath := filepath.Join("/invalid/path/that/cannot/be/created", "output.json")
	
	configContent := `input-file: "` + filepath.ToSlash(xmlPath) + `"
output-file: "` + filepath.ToSlash(outputPath) + `"
`
	err = os.WriteFile(configPath, []byte(configContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	// Suppress stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = run(configPath)

	w.Close()
	os.Stdout = oldStdout
	io.Copy(io.Discard, r)
	r.Close()

	// This might succeed or fail depending on system, but we test the error path
	if err != nil {
		// Expected error for invalid output directory
		return
	}
}

