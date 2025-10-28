package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Svyatoslav2324/task-3/internal/marshaljson"
	"github.com/Svyatoslav2324/task-3/internal/structures"
	"github.com/Svyatoslav2324/task-3/internal/unmarshalxml"
	"github.com/Svyatoslav2324/task-3/internal/unmarshalyaml"
)

func makeDir(outFile string) {
	dir := filepath.Dir(outFile)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0o755)
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

	inputFile, err := os.Open(config.InputFile)
	if err != nil {
		panic(err)
	}
	defer closeFile(inputFile)

	valutes := new(structures.XMLStruct)

	err = unmarshalxml.UnMarshalXML(inputFile, valutes)
	if err != nil {
		fmt.Println(err)

		return
	}

	sort.Slice(valutes.ValCurs, func(i, j int) bool {
		valuei, _ := strconv.ParseFloat(strings.Replace(valutes.ValCurs[i].Value, ",", ".", 1), 64)
		valuej, _ := strconv.ParseFloat(strings.Replace(valutes.ValCurs[j].Value, ",", ".", 1), 64)

		return valuei > valuej
	})

	makeDir(config.OutputFile)

	outputFile, _ := os.Create(config.OutputFile)
	defer closeFile(outputFile)

	jsonValutes, err := structures.ConvertXMLToJSON(*valutes)
	if err != nil {
		fmt.Println(err)

		return
	}

	err = marshaljson.MarshalJSON(outputFile, &jsonValutes)
	if err != nil {
		fmt.Println(err)

		return
	}
}
