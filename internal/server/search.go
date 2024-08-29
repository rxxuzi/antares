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

type SearchCondition struct {
	Term    string
	IsRegex bool
	IsNot   bool
}

type SearchQuery struct {
	Conditions []SearchCondition
	IsOr       bool
}

func parseQuery(query string, useRegex bool) SearchQuery {
	terms := strings.Fields(query)
	var sq SearchQuery
	sq.IsOr = false

	for i := 0; i < len(terms); i++ {
		term := terms[i]

		if strings.EqualFold(term, "OR") || term == "||" {
			sq.IsOr = true
			continue
		}

		if strings.EqualFold(term, "AND") || term == "&&" {
			continue
		}

		condition := SearchCondition{
			Term:    term,
			IsRegex: useRegex,
			IsNot:   false,
		}

		if strings.HasPrefix(term, "NOT") || strings.HasPrefix(term, "!!") {
			condition.IsNot = true
			condition.Term = strings.TrimPrefix(strings.TrimPrefix(term, "NOT"), "!!")
		}

		sq.Conditions = append(sq.Conditions, condition)
	}
	return sq
}

func (sq SearchQuery) match(filename string, caseSensitive bool) bool {
	matchCount := 0

	for _, condition := range sq.Conditions {
		matches := false
		if condition.IsRegex {
			re, err := regexp.Compile(condition.Term)
			if err == nil {
				matches = re.MatchString(filename)
			}
		} else {
			if !caseSensitive {
				filename = strings.ToLower(filename)
				condition.Term = strings.ToLower(condition.Term)
			}
			matches = strings.Contains(filename, condition.Term)
		}

		if condition.IsNot {
			matches = !matches
		}

		if matches {
			matchCount++
		}

		if sq.IsOr && matches {
			return true
		}
	}

	return !sq.IsOr && matchCount == len(sq.Conditions)
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

	searchQuery := parseQuery(query, useRegex)
	results, err := searchFiles(root, searchQuery, caseSensitive)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error searching files: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received query: %s, Regex: %v, Case Sensitive: %v\n", query, useRegex, caseSensitive)

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

func searchFiles(root string, query SearchQuery, caseSensitive bool) ([]SearchResult, error) {
	var results []SearchResult

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsPermission(err) {
				// Skip directories with permission errors
				fmt.Printf("Skipping inaccessible path: %s\n", path)
				return filepath.SkipDir
			}
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		if query.match(info.Name(), caseSensitive) {
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

	if err != nil && !os.IsPermission(err) {
		return nil, err
	}

	fmt.Printf("Query: %+v, Results count: %d\n", query, len(results))

	return results, nil
}
