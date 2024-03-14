package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/image"
	"github.com/dunkbing/tinyimg/converter/utils"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RequestBody struct {
	Files []string `json:"files"`
}

func getContentType(fileName string) string {
	switch filepath.Ext(fileName) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

func isImage(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	return strings.HasPrefix(mimeType, "image/")
}

// isFileUploaded checks if a file with the given name has already been uploaded.
func isFileUploaded(filePath string) bool {
	// Check if the file exists in the directory
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

type handler struct {
	fileManager *image.FileManager
	config      *config.Config
}

func New() *handler {
	c := config.GetConfig()
	return &handler{
		fileManager: image.NewFileManager(),
		config:      c,
	}
}

func (h *handler) Upload(w http.ResponseWriter, r *http.Request) {
	var sizeLimit int64 = 10 * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, sizeLimit)

	startTime := time.Now()
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file. The file may be too large (max 10MB)", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	mimeType := http.DetectContentType(data)
	if !isImage(mimeType) {
		http.Error(w, "Invalid file format. Only images are allowed.", http.StatusBadRequest)
		return
	}
	fileType, _ := image.GetFileType(mimeType)
	filename, err := utils.GenerateHash(data)
	if err != nil {
		slog.Error("Error generating file name", "err", err.Error())
		filename = filepath.Base(header.Filename)
	}
	ext := fmt.Sprintf(".%s", fileType)
	filename = fmt.Sprintf("%s%s", filename, ext)

	dest := filepath.Join(h.config.App.InDir, filename)
	if !isFileUploaded(dest) {
		slog.Info("Upload", "dest", dest)
		err = os.WriteFile(dest, data, 0644)
		if err != nil {
			http.Error(w, "Error writing the file", http.StatusInternalServerError)
			return
		}
	}
	took := time.Since(startTime).Seconds()
	fmt.Println("Write to file took", took, "seconds")

	formatStr := r.FormValue("formats")
	formats := make([]string, 0)
	if formatStr != "" {
		formats = strings.Split(formatStr, ",")
	} else {
		formats = append(formats, fileType)
	}

	f := image.File{
		Data:          data,
		Ext:           ext,
		MimeType:      mimeType,
		Name:          filename,
		Size:          header.Size,
		Formats:       formats,
		InputFileDest: dest,
	}

	err = h.fileManager.HandleFile(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	results, files, errs := h.fileManager.Convert()
	strErrs := make([]string, len(errs))
	for i, err := range errs {
		strErrs[i] = err.Error()
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"data":   results,
		"files":  files,
		"errors": strErrs,
	})
}

func (h *handler) DownloadAll(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	zippedPath, err := h.fileManager.ZipFiles(body.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	zipFileName := filepath.Base(zippedPath)
	zipFile, err := os.Open(zippedPath)
	if err != nil {
		http.Error(w, "Error opening zip file", http.StatusInternalServerError)
		return
	}
	defer zipFile.Close()

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	_, err = io.Copy(w, zipFile)
	if err != nil {
		http.Error(w, "Error serving zip file", http.StatusInternalServerError)
		return
	}
}

func (h *handler) ServeImg(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("f")
	if fileName == "" {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(h.config.App.OutDir, fileName)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	contentType := getContentType(fileName)
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}
