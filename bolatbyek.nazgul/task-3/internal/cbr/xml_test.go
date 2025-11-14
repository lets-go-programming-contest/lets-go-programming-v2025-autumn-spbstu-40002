package cbr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseXML(t *testing.T) {
	// Create a temporary XML file
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "test.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<ValCurs Date="02.03.2022" name="Foreign Currency Market">
  <Valute ID="R01035">
    <NumCode>826</NumCode>
    <CharCode>GBP</CharCode>
    <Nominal>1</Nominal>
    <Name>Фунт стерлингов Соединенного королевства</Name>
    <Value>95,0000</Value>
  </Valute>
  <Valute ID="R01060">
    <NumCode>051</NumCode>
    <CharCode>AMD</CharCode>
    <Nominal>100</Nominal>
    <Name>Армянских драмов</Name>
    <Value>12,5000</Value>
  </Valute>
</ValCurs>
`
	err := os.WriteFile(xmlPath, []byte(xmlContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	valCurs, err := ParseXML(xmlPath)
	if err != nil {
		t.Fatalf("ParseXML failed: %v", err)
	}

	if valCurs == nil {
		t.Fatal("ParseXML returned nil ValCurs")
	}

	if len(valCurs.Valutes) != 2 {
		t.Errorf("Expected 2 valutes, got %d", len(valCurs.Valutes))
	}

	if valCurs.Valutes[0].CharCode != "GBP" {
		t.Errorf("Expected first CharCode to be 'GBP', got '%s'", valCurs.Valutes[0].CharCode)
	}

	if valCurs.Valutes[0].NumCode != "826" {
		t.Errorf("Expected first NumCode to be '826', got '%s'", valCurs.Valutes[0].NumCode)
	}

	if valCurs.Valutes[0].Value != "95,0000" {
		t.Errorf("Expected first Value to be '95,0000', got '%s'", valCurs.Valutes[0].Value)
	}
}

func TestParseXML_NonExistentFile(t *testing.T) {
	_, err := ParseXML("non_existent_file.xml")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestParseXML_InvalidXML(t *testing.T) {
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "invalid.xml")
	xmlContent := `invalid xml content`
	err := os.WriteFile(xmlPath, []byte(xmlContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	_, err = ParseXML(xmlPath)
	if err == nil {
		t.Error("Expected error for invalid XML, got nil")
	}
}

func TestParseXML_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	xmlPath := filepath.Join(tmpDir, "empty.xml")
	err := os.WriteFile(xmlPath, []byte(""), 0600)
	if err != nil {
		t.Fatalf("Failed to create test XML file: %v", err)
	}

	_, err = ParseXML(xmlPath)
	if err == nil {
		t.Error("Expected error for empty XML file, got nil")
	}
}

