package rawConfig

import (
	"embed"
)

//go:embed "config.json"
var TPV2ConfigFs embed.FS

var Version string
