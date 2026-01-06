//go:build dev

package config

import (
	"embed"
)

//go:embed dev.yaml
var confFile embed.FS

var filePath = "dev.yaml"

func getFilePath() string {
	return filePath
}
