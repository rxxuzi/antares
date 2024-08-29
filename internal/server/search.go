package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
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
	query, err := url.QueryUnescape(r.URL.Query().Get("q"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid query: %v", err), http.StatusBadRequest)
		return
	}
	useRegex := r.URL.Query().Get("r") == "true" || r.URL.Query().Has("r")
	caseSensitive := r.URL.Query().Get("c") == "true" || r.URL.Query().Has("c")

	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received query: %s, Regex: %v, Case Sensitive: %v\n", query, useRegex, caseSensitive)

	results, err := searchFiles(root, query, caseSensitive, useRegex)
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

func searchFiles(root, query string, caseSensitive, useRegex bool) ([]SearchResult, error) {
	var results []SearchResult
	var matcher func(string) bool

	if useRegex {
		var regexFlags string
		if !caseSensitive {
			regexFlags = "(?i)"
		}
		// Escape the query if it's not already a valid regex
		if _, err := regexp.Compile(query); err != nil {
			query = regexp.QuoteMeta(query)
		}
		re, err := regexp.Compile(regexFlags + query)
		if err != nil {
			return nil, fmt.Errorf("invalid regex: %v", err)
		}
		matcher = re.MatchString
	} else {
		if !caseSensitive {
			query = strings.ToLower(query)
		}
		matcher = func(s string) bool {
			if !caseSensitive {
				s = strings.ToLower(s)
			}
			return strings.Contains(s, query)
		}
	}

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

		if matcher(info.Name()) {
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
