//go:build !dev

package config

import (
	"embed"
)

//go:embed prod.yaml
var confFile embed.FS

var filePath = "prod.yaml"

func getFilePath() string {
	return filePath
}
