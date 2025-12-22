//go:build !prod

package config

import (
    "fmt"

    _ "embed"
    "gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devYaml []byte

func Load() (Config, error) {
    var cfg Config

    err := yaml.Unmarshal(devYaml, &cfg)
    if err != nil {
        return Config{}, fmt.Errorf("failed to unmarshal dev config: %w", err)
    }

    return cfg, nil
}
