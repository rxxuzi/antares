package server

import (
	"fmt"
	"github.com/rxxuzi/antares/internal/static"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func antaresServer(root string, staticFS fs.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(root, filepath.Clean(r.URL.Path))

		if r.URL.Path == "/search" {
			searchHandler(w, r, root, staticFS)
			return
		}

		if r.Method == "POST" {
			handleFileUpload(w, r, root)
			return
		}

		info, err := os.Stat(path)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		if !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}

		files, err := os.ReadDir(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to read directory: %v", err), http.StatusInternalServerError)
			return
		}

		var fileInfos []FileInfo
		for _, file := range files {
			info, _ := file.Info()
			fileType := "folder"
			if !file.IsDir() {
				fileType = GetFileType(file.Name(), path)
			}
			fileInfos = append(fileInfos, FileInfo{
				Name:    file.Name(),
				Size:    info.Size(),
				Mode:    info.Mode(),
				ModTime: info.ModTime(),
				IsDir:   file.IsDir(),
				Type:    fileType,
			})
		}

		sort.Slice(fileInfos, func(i, j int) bool {
			if fileInfos[i].IsDir != fileInfos[j].IsDir {
				return fileInfos[i].IsDir
			}
			return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
		})

		tmplContent, err := fs.ReadFile(staticFS, "antares.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read template: %v", err), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.New("listing").Funcs(template.FuncMap{
			"formatSize": formatSize,
			"formatTime": formatTime,
			"isImage":    isImage,
			"isVideo":    isVideo,
			"truncate":   truncate,
			"getFileIcon": func(name string) string {
				for _, file := range fileInfos {
					if file.Name == name {
						return GetFileIcon(name, file.Type)
					}
				}
				return "fas fa-file"
			},
			"split": strings.Split,
		}).Parse(string(tmplContent))

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse template: %v", err), http.StatusInternalServerError)
			return
		}

		data := DirectoryContent{
			Path:        r.URL.Path,
			Files:       fileInfos,
			CurrentPath: filepath.Base(path),
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, fmt.Sprintf("Failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func handleFileUpload(w http.ResponseWriter, r *http.Request, root string) {
	// 32 MB のメモリ制限を設定
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	safeFileName := filepath.Base(handler.Filename)
	uploadPath := filepath.Join(root, r.URL.Path, safeFileName)

	for i := 1; fileExists(uploadPath); i++ {
		ext := filepath.Ext(safeFileName)
		name := strings.TrimSuffix(safeFileName, ext)
		uploadPath = filepath.Join(root, r.URL.Path, fmt.Sprintf("%s_%d%s", name, i, ext))
	}

	f, err := os.Create(uploadPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}

// rootHandler handles the root path and redirects to /drive/
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, PREFIX_DRIVE, http.StatusFound)
		return
	}
	custom404Handler(w)
}

func custom404Handler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	staticFS := static.GetFS()
	custom404, err := fs.ReadFile(staticFS, "404.html")
	if err != nil {
		http.Error(w, "<html><head><title>404 Not Found</title></head><body><h1 style=\"text-align: center\">404 Not Found</h1></body></html>", http.StatusNotFound)
		return
	}
	w.Write(custom404)
}
