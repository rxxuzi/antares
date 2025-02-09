package static

import (
	"embed"
	"io/fs"
)

//go:embed web
var WebFS embed.FS

//go:embed doc
var DocFS embed.FS

func GetFS() fs.FS {
	staticFS, _ := fs.Sub(WebFS, "web")
	return staticFS
}

func GetDoc() fs.FS {
	docFS, _ := fs.Sub(DocFS, "doc")
	return docFS
}
