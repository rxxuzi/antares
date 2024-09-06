package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type APIRequest struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func handleAPI(w http.ResponseWriter, r *http.Request, root string) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req APIRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONResponse(w, false, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	switch req.Type {
	case "rm":
		handleDelete(w, req, root)
	default:
		sendJSONResponse(w, false, "Unknown operation type", http.StatusBadRequest)
	}
}

func handleDelete(w http.ResponseWriter, req APIRequest, root string) {
	if req.Path == "" {
		sendJSONResponse(w, false, "Path is required", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(root, filepath.Clean(req.Path))

	// Ensure the path is within the root directory
	if !strings.HasPrefix(fullPath, root) {
		sendJSONResponse(w, false, "Invalid file path", http.StatusBadRequest)
		return
	}

	// Check if file exists
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		sendJSONResponse(w, false, "File not found", http.StatusNotFound)
		return
	}

	// Delete the file
	err = os.Remove(fullPath)
	if err != nil {
		log.Printf("Failed to delete file: %v", err)
		sendJSONResponse(w, false, fmt.Sprintf("Failed to delete file: %v", err), http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, true, fmt.Sprintf("File deleted successfully: %s", req.Path), http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, success bool, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Success: success,
		Message: message,
	})
}
