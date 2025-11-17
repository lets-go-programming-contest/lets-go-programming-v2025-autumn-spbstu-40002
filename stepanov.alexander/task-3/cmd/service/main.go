package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

type Currency struct {
	NumCode  int     `json:"num_code" xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value" xml:"Value"`
}

type CurrencyList struct {
	Currencies []Currency `xml:"Valute"`
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}

func parseValue(valueStr string) (float64, error) {
	valueStr = strings.Replace(valueStr, ",", ".", 1)
	return strconv.ParseFloat(valueStr, 64)
}

func (c *Currency) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	type currencyXML struct {
		NumCode  string `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Value    string `xml:"Value"`
	}

	var temp currencyXML
	if err := decoder.DecodeElement(&temp, &start); err != nil {
		return err
	}

	numCode, err := strconv.Atoi(temp.NumCode)
	if err != nil {
		return fmt.Errorf("parse num code: %w", err)
	}

	value, err := parseValue(temp.Value)
	if err != nil {
		return fmt.Errorf("parse value: %w", err)
	}

	c.NumCode = numCode
	c.CharCode = temp.CharCode
	c.Value = value

	return nil
}

func loadCurrencies(path string) ([]Currency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if strings.ToLower(charset) == "windows-1251" {
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		}
		return input, nil
	}

	var data struct {
		Currencies []Currency `xml:"Valute"`
	}

	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("decode xml: %w", err)
	}

	return data.Currencies, nil
}

func sortCurrencies(currencies []Currency) {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})
}

func saveCurrencies(path string, currencies []Currency) error {
	if err := os.MkdirAll(path[:strings.LastIndex(path, "/")], 0755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(currencies); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("config file path is required")
	}

	cfg, err := loadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	currencies, err := loadCurrencies(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	sortCurrencies(currencies)

	if err := saveCurrencies(cfg.OutputFile, currencies); err != nil {
		panic(err)
	}

	fmt.Println("Successfully processed currencies")
}
