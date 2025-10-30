package unmarshalyaml

import (
	"flag"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func GetPaths() Config {
	configPath := flag.String("config", "", "path to YAML config")
	flag.Parse()

	if *configPath == "" {
		panic("Empty config path")
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}

	return config
}
