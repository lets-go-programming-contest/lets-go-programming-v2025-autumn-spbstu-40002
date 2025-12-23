//go:build !dev && !prod

package config

func LoadConfig() (*Config, error) {
	return loadConfig(defaultConfig)
}
