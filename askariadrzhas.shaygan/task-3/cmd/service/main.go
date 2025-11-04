package main

import (
	"fmt"

	"github.com/XShaygaND/task-3/internal/config"
	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	settings := config.LoadSettings()
	currencyData := parser.ExtractCurrencyData(settings.SourcePath)
	sortedData := processor.OrganizeByRate(currencyData)
	writer.SaveAsJSON(sortedData, settings.TargetPath)

	fmt.Printf("Processed %d currency records\n", len(sortedData))
	fmt.Printf("Output saved to: %s\n", settings.TargetPath)
}
