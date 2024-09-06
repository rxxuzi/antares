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
	File   bool                   `json:"file"`
	Type   string                 `json:"type"`
	Path   string                 `json:"path"`
	Src    string                 `json:"src,omitempty"`
	Dst    string                 `json:"dst,omitempty"`
	Option map[string]interface{} `json:"option,omitempty"`
}

type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

const (
	ErrInvalidMethod    = "INVALID_METHOD"
	ErrInvalidJSON      = "INVALID_JSON"
	ErrUnknownOperation = "UNKNOWN_OPERATION"
	ErrMissingPath      = "MISSING_PATH"
	ErrInvalidPath      = "INVALID_PATH"
	ErrFileNotFound     = "FILE_NOT_FOUND"
	ErrOperationFailed  = "OPERATION_FAILED"
)

func handleAPI(w http.ResponseWriter, r *http.Request, root string) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, false, "Method not allowed", ErrInvalidMethod, http.StatusMethodNotAllowed)
		return
	}

	var req APIRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sendJSONResponse(w, false, "Invalid JSON payload", ErrInvalidJSON, http.StatusBadRequest)
		return
	}

	switch req.Type {
	case "delete":
		handleDelete(w, req, root)
	case "move":
		handleMove(w, req, root)
	case "copy":
		handleCopy(w, req, root)
	case "rename":
		handleRename(w, req, root)
	default:
		sendJSONResponse(w, false, "Unknown operation type", ErrUnknownOperation, http.StatusBadRequest)
	}
}

func handleDelete(w http.ResponseWriter, req APIRequest, root string) {
	if req.Path == "" {
		sendJSONResponse(w, false, "Path is required", ErrMissingPath, http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(root, filepath.Clean(req.Path))

	if !strings.HasPrefix(fullPath, root) {
		sendJSONResponse(w, false, "Invalid file path", ErrInvalidPath, http.StatusBadRequest)
		return
	}

	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		sendJSONResponse(w, false, "File or folder not found", ErrFileNotFound, http.StatusNotFound)
		return
	}

	var deleteErr error
	if req.File {
		deleteErr = os.Remove(fullPath)
	} else {
		deleteErr = os.RemoveAll(fullPath)
	}

	if deleteErr != nil {
		fileOrFolder := "file"
		if !req.File {
			fileOrFolder = "folder"
		}
		log.Printf("Failed to delete %s: %v", fileOrFolder, deleteErr)
		sendJSONResponse(w, false, fmt.Sprintf("Failed to delete %s: %v", fileOrFolder, deleteErr), ErrOperationFailed, http.StatusInternalServerError)
		return
	}

	fileOrFolder := "File"
	if !req.File {
		fileOrFolder = "Folder"
	}
	sendJSONResponse(w, true, fmt.Sprintf("%s deleted successfully: %s", fileOrFolder, req.Path), "", http.StatusOK)
}

func handleMove(w http.ResponseWriter, req APIRequest, root string) {
	// TODO: Implement
	sendJSONResponse(w, false, "Move operation not implemented yet", ErrOperationFailed, http.StatusNotImplemented)
}

func handleCopy(w http.ResponseWriter, req APIRequest, root string) {
	// TODO: Implement
	sendJSONResponse(w, false, "Copy operation not implemented yet", ErrOperationFailed, http.StatusNotImplemented)
}

func handleRename(w http.ResponseWriter, req APIRequest, root string) {
	if req.Path == "" || req.Dst == "" {
		sendJSONResponse(w, false, "Both source and destination paths are required", ErrMissingPath, http.StatusBadRequest)
		return
	}

	srcPath := filepath.Join(root, filepath.Clean(req.Path))
	dstPath := filepath.Join(root, filepath.Clean(req.Dst))

	log.Printf("src -> %s\n", srcPath)
	log.Printf("dst -> %s\n", dstPath)

	if !strings.HasPrefix(srcPath, root) || !strings.HasPrefix(dstPath, root) {
		sendJSONResponse(w, false, "Invalid file path", ErrInvalidPath, http.StatusBadRequest)
		return
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		sendJSONResponse(w, false, "Source file or folder not found", ErrFileNotFound, http.StatusNotFound)
		return
	}

	if err := os.Rename(srcPath, dstPath); err != nil {
		log.Printf("Failed to rename: %v", err)
		sendJSONResponse(w, false, fmt.Sprintf("Failed to rename: %v", err), ErrOperationFailed, http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, true, fmt.Sprintf("Successfully renamed from %s to %s", req.Path, req.Dst), "", http.StatusOK)
}

func sendJSONResponse(w http.ResponseWriter, success bool, message string, errorCode string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(APIResponse{
		Success:   success,
		Message:   message,
		ErrorCode: errorCode,
	})
}
