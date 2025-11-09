package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// Config структура для конфигурационного файла YAML
type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

// ValCurs структура для XML данных с сайта ЦБ РФ
type ValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valutes []Valute `xml:"Valute"`
}

// Valute структура для валюты
type Valute struct {
	XMLName  xml.Name `xml:"Valute"`
	NumCode  string   `xml:"NumCode"`
	CharCode string   `xml:"CharCode"`
	Value    string   `xml:"Value"`
}

// CurrencyResult структура для результата в JSON
type CurrencyResult struct {
	NumCode  string  `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func main() {
	// Шаг 1: Чтение конфигурации
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic("Config path is required")
	}

	config, err := loadConfig(*configPath)
	if err != nil {
		panic(fmt.Sprintf("Error loading config: %v", err))
	}

	// Шаг 2: Чтение и декодирование XML
	valCurs, err := parseXML(config.InputFile)
	if err != nil {
		panic(fmt.Sprintf("Error parsing XML: %v", err))
	}

	// Шаг 3: Преобразование и сортировка данных
	currencies := convertAndSortCurrencies(valCurs)

	// Шаг 4: Сохранение результатов в JSON
	err = saveResults(config.OutputFile, currencies)
	if err != nil {
		panic(fmt.Sprintf("Error saving results: %v", err))
	}

	fmt.Println("Successfully processed currencies")
}

// loadConfig загружает конфигурацию из YAML файла
func loadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// parseXML парсит XML файл с данными о валютах
func parseXML(filePath string) (*ValCurs, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var valCurs ValCurs
	err = xml.Unmarshal(data, &valCurs)
	if err != nil {
		return nil, err
	}

	return &valCurs, nil
}

// convertAndSortCurrencies конвертирует и сортирует валюты по убыванию значения
func convertAndSortCurrencies(valCurs *ValCurs) []CurrencyResult {
	var currencies []CurrencyResult

	for _, valute := range valCurs.Valutes {
		// Преобразуем значение из строки в float64
		var value float64
		_, err := fmt.Sscanf(valute.Value, "%f", &value)
		if err != nil {
			// Пропускаем валюты с некорректными значениями
			continue
		}

		currencies = append(currencies, CurrencyResult{
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Value:    value,
		})
	}

	// Сортируем по убыванию значения
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	return currencies
}

// saveResults сохраняет результаты в JSON файл
func saveResults(outputPath string, currencies []CurrencyResult) error {
	// Создаем директорию если она не существует
	dir := filepath.Dir(outputPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")

	return encoder.Encode(currencies)
}
