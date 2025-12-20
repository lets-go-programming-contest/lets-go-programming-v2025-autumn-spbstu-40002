package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Svyatoslav2324/task-3/internal/marshaljson"
	"github.com/Svyatoslav2324/task-3/internal/structures"
	"github.com/Svyatoslav2324/task-3/internal/unmarshalxml"
	"github.com/Svyatoslav2324/task-3/internal/unmarshalyaml"
)

func makeDir(outFile string) {
	dir := filepath.Dir(outFile)
	dirPerm := 0o755

	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.FileMode(dirPerm))
		if err != nil {
			panic(err)
		}
	}
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	config := unmarshalyaml.GetPaths()

	makeDir(config.OutputFile)

	outputFile, _ := os.Create(config.OutputFile)
	defer closeFile(outputFile)

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		panic(err)
	}
	defer closeFile(inputFile)

	content, err := io.ReadAll(inputFile)
	if err != nil {
		panic(err)
	}

	newContent := strings.ReplaceAll(string(content), ",", ".")

	_, err = outputFile.WriteString(newContent)
	if err != nil {
		panic(err)
	}

	_, err = outputFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	valutes := new(structures.ValuteStruct)

	err = unmarshalxml.UnMarshalXML(outputFile, valutes)
	if err != nil {
		panic(err)
	}

	_, err = outputFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	err = outputFile.Truncate(0)
	if err != nil {
		panic(err)
	}

	sort.Slice(valutes.ValCurs, func(i, j int) bool {
		return valutes.ValCurs[i].Value > valutes.ValCurs[j].Value
	})

	err = marshaljson.MarshalJSON(outputFile, valutes)
	if err != nil {
		fmt.Println(err)

		return
	}
}
