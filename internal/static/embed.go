package static

import (
	"embed"
	"io/fs"
)

//go:embed web
var WebFS embed.FS

func GetFS() fs.FS {
	staticFS, _ := fs.Sub(WebFS, "web")
	return staticFS
}
