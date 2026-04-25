package embedfs

import "embed"

//go:embed all:web/dist
var Static embed.FS
