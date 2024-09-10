package server

import (
	"fmt"
	"github.com/rxxuzi/antares/internal/static"
	"net/http"
)

const (
	PREFIX_DRIVE string = "/drive/"
)

type DirectoryContent struct {
	Path        string
	Files       []FileInfo
	CurrentPath string
}

// CreateServer creates an HTTP server with the specified configuration
func CreateServer(config *Config) (*http.Server, error) {
	mux := http.NewServeMux()
	staticFS := static.GetFS()
	var handler http.Handler = mux

	mux.HandleFunc("/", rootHandler)

	// 静的ファイル
	fileServer := http.FileServer(http.FS(staticFS))
	mux.Handle("/web/", http.StripPrefix("/web/", fileServer))

	// 検索ハンドラー
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		searchHandler(w, r, config.RootDir, staticFS)
	})

	antares := antaresServer(config.RootDir, staticFS)
	mux.Handle(PREFIX_DRIVE, http.StripPrefix(PREFIX_DRIVE, antares))

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		handleAPI(w, r, config.RootDir)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r, staticFS)
	})

	mux.HandleFunc("/ws", wsHandler)

	// ログ記録用のミドルウェア）
	if config.LogFlag {
		handler = logRequest(handler)
	}

	// 利用可能なポートを見つける
	availablePort, err := findAvailablePort(config.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to find available port: %v", err)
	}

	config.Port = availablePort

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: handler,
	}, nil
}
