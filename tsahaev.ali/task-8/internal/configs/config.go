package configs

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

type ConfigType struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

func PrintTagAndLevel() error {
	var conf ConfigType

	err := yaml.Unmarshal(dataYaml, &conf)
	if err != nil {
		return fmt.Errorf("error unmarshal: %w", err)
	}

	fmt.Print(conf.Environment + " " + conf.LogLevel)

	return nil
}
