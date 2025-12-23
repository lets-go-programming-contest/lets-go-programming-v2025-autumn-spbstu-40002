package configs

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v2"
)

type Settings struct {
	Env string `yaml:"environment"`
	Log string `yaml:"log_level"`
}

func DisplayEnvironmentAndLogLevel() error {
	var cfg Settings

	err := yaml.Unmarshal(configData, &cfg)
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	fmt.Print(cfg.Env + " " + cfg.Log)

	return nil
}
