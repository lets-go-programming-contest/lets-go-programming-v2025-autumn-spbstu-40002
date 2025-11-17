package convert

import (
	"encoding/xml"
	"testing"

	"github.com/Nazkaaa/task-3/internal/cbr"
)

func TestConvertAndSort(t *testing.T) {
	valCurs := &cbr.ValCurs{
		XMLName: xml.Name{Local: "ValCurs"},
		Valutes: []cbr.Valute{
			{
				NumCode:  "826",
				CharCode: "GBP",
				Value:    "95,0000",
			},
			{
				NumCode:  "840",
				CharCode: "USD",
				Value:    "120,0000",
			},
			{
				NumCode:  "978",
				CharCode: "EUR",
				Value:    "130,0000",
			},
		},
	}

	result := ConvertAndSort(valCurs)

	if len(result) != 3 {
		t.Errorf("Expected 3 currencies, got %d", len(result))
	}

	// Check sorting (should be descending by value)
	if result[0].Value != 130.0 {
		t.Errorf("Expected first value to be 130.0, got %f", result[0].Value)
	}
	if result[0].CharCode != "EUR" {
		t.Errorf("Expected first CharCode to be 'EUR', got '%s'", result[0].CharCode)
	}

	if result[1].Value != 120.0 {
		t.Errorf("Expected second value to be 120.0, got %f", result[1].Value)
	}
	if result[1].CharCode != "USD" {
		t.Errorf("Expected second CharCode to be 'USD', got '%s'", result[1].CharCode)
	}

	if result[2].Value != 95.0 {
		t.Errorf("Expected third value to be 95.0, got %f", result[2].Value)
	}
	if result[2].CharCode != "GBP" {
		t.Errorf("Expected third CharCode to be 'GBP', got '%s'", result[2].CharCode)
	}

	// Check NumCode conversion
	if result[0].NumCode != 978 {
		t.Errorf("Expected first NumCode to be 978, got %d", result[0].NumCode)
	}
}

func TestConvertAndSort_EmptyList(t *testing.T) {
	valCurs := &cbr.ValCurs{
		XMLName: xml.Name{Local: "ValCurs"},
		Valutes: []cbr.Valute{},
	}

	result := ConvertAndSort(valCurs)

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d items", len(result))
	}
}

func TestConvertAndSort_SingleItem(t *testing.T) {
	valCurs := &cbr.ValCurs{
		XMLName: xml.Name{Local: "ValCurs"},
		Valutes: []cbr.Valute{
			{
				NumCode:  "840",
				CharCode: "USD",
				Value:    "120,0000",
			},
		},
	}

	result := ConvertAndSort(valCurs)

	if len(result) != 1 {
		t.Errorf("Expected 1 currency, got %d", len(result))
	}

	if result[0].CharCode != "USD" {
		t.Errorf("Expected CharCode to be 'USD', got '%s'", result[0].CharCode)
	}

	if result[0].Value != 120.0 {
		t.Errorf("Expected Value to be 120.0, got %f", result[0].Value)
	}
}