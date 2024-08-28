package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type SearchResult struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Type    string    `json:"type"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"modTime"`
}

type SearchResultPage struct {
	Query   string
	Results []SearchResult
}

func searchHandler(w http.ResponseWriter, r *http.Request, root string, staticFS fs.FS) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	results, err := searchFiles(root, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching files: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the client wants JSON (e.g., for AJAX requests)
	if r.Header.Get("Accept") == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
		return
	}

	// Render HTML page
	tmplContent, err := fs.ReadFile(staticFS, "search-result.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read template: %v", err), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("search-result").Funcs(template.FuncMap{
		"truncate":    truncate,
		"getFileIcon": GetFileIcon,
		"isImage":     isImage,
		"isVideo":     isVideo,
		"formatSize":  formatSize,
		"formatTime":  formatTime,
	}).Parse(string(tmplContent))

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
		return
	}

	data := SearchResultPage{
		Query:   query,
		Results: results,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
		return
	}
}

func searchFiles(root, query string) ([]SearchResult, error) {
	var results []SearchResult

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(info.Name()), strings.ToLower(query)) {
			fileType := GetFileType(info.Name(), root)
			results = append(results, SearchResult{
				Name:    info.Name(),
				Path:    filepath.ToSlash(relPath),
				Type:    fileType,
				Size:    info.Size(),
				ModTime: info.ModTime(),
			})
		}

		return nil
	})

	return results, err
}
