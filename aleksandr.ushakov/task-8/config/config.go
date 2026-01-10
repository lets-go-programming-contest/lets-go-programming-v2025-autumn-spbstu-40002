package config

type Config struct {
	Environment string `yaml:"environment"`
	Loglevel    string `yaml:"log_level"`
}
