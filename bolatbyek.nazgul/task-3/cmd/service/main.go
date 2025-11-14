package main

import (
    "flag"
    "fmt"

    "github.com/Nazkaaa/task-3/internal/cbr"
    "github.com/Nazkaaa/task-3/internal/config"
    "github.com/Nazkaaa/task-3/internal/convert"
    "github.com/Nazkaaa/task-3/internal/output"
)

func main() {
    configPath := flag.String("config", "config.yaml", "Path to the YAML configuration file")
    flag.Parse()

    cfg, err := config.LoadConfig(*configPath)
    if err != nil {
        panic("Error loading config: " + err.Error())
    }

    err = config.EnsureOutputDir(cfg.OutputFile)
    if err != nil {
        panic("Error creating output directory: " + err.Error())
    }

    valCurs, err := cbr.ParseXML(cfg.InputFile)
    if err != nil {
        panic("Error parsing XML: " + err.Error())
    }

    currencies := convert.ConvertAndSort(valCurs)

    var outputCurrencies []interface{}
    for _, c := range currencies {
        outputCurrencies = append(outputCurrencies, c)
    }

    err = output.SaveToJSON(outputCurrencies, cfg.OutputFile)
    if err != nil {
        panic("Error saving to JSON: " + err.Error())
    }

    fmt.Printf("Successfully processed %d currencies and saved to %s\n", len(currencies), cfg.OutputFile)
}
