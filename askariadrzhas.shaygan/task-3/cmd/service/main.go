package main

import (
	"fmt"

	"github.com/XShaygaND/task-3/internal/config"
	"github.com/XShaygaND/task-3/internal/parser"
	"github.com/XShaygaND/task-3/internal/processor"
	"github.com/XShaygaND/task-3/internal/writer"
)

func main() {
	settings, err := config.LoadSettings()
	if err != nil {
		return
	}

	data := parser.ExtractCurrencyData(settings.SourcePath)
	sorted := processor.OrganizeByRate(data)
	writer.SaveAsJSON(sorted, settings.TargetPath)

	fmt.Printf("Processed %d records\n", len(sorted))
}
