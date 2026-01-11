package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveToJSON(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_output.json")

	currencies := []interface{}{
		map[string]interface{}{
			"num_code":  840,
			"char_code": "USD",
			"value":     120.0,
		},
		map[string]interface{}{
			"num_code":  978,
			"char_code": "EUR",
			"value":     130.0,
		},
	}

	err := SaveToJSON(currencies, outputPath)
	if err != nil {
		t.Fatalf("SaveToJSON failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %v", err)
	}

	// Read and verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result []interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 currencies in JSON, got %d", len(result))
	}
}

func TestSaveToJSON_EmptyList(t *testing.T) {
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "empty_output.json")

	currencies := []interface{}{}

	err := SaveToJSON(currencies, outputPath)
	if err != nil {
		t.Fatalf("SaveToJSON failed for empty list: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %v", err)
	}

	// Read and verify content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	var result []interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty array in JSON, got %d items", len(result))
	}
}

func TestSaveToJSON_WithExistingDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "subdir", "nested")
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	outputPath := filepath.Join(outputDir, "output.json")

	currencies := []interface{}{
		map[string]interface{}{
			"num_code":  840,
			"char_code": "USD",
			"value":     120.0,
		},
	}

	err = SaveToJSON(currencies, outputPath)
	if err != nil {
		t.Fatalf("SaveToJSON failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Fatalf("Output file was not created: %v", err)
	}
}

