package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"optipic/converter/image"
	"path/filepath"
	"strings"
)

type RequestBody struct {
	Files []string `json:"files"`
}

var allowedOrigin = "*"

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	mimeType := header.Header.Get("Content-Type")
	fileType, _ := image.GetFileType(mimeType)
	filename := header.Filename
	ext := filepath.Ext(header.Filename)
	formatStr := r.FormValue("formats")
	formats := make([]string, 0)
	if formatStr != "" {
		formats = strings.Split(formatStr, ",")
	} else {
		formats = append(formats, fileType)
	}

	f := image.File{
		Data:     data,
		Ext:      ext,
		MimeType: header.Header.Get("Content-Type"),
		Name:     filename,
		Size:     header.Size,
		Formats:  formats,
	}

	if !isImage(f.MimeType) {
		http.Error(w, "Invalid file format. Only images are allowed.", http.StatusBadRequest)
		return
	}

	fileManager := image.NewFileManager()
	err = fileManager.HandleFile(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	results, files, errs := fileManager.Convert()
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

func downloadZipHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	fm := image.NewFileManager()
	zippedUrl, err := fm.ZipFiles(body.Files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"url": zippedUrl,
	})
}

func isImage(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	return strings.HasPrefix(mimeType, "image/")
}

func main() {
	mux := http.NewServeMux()
	uploadHandlerWithCors := enableCors(http.HandlerFunc(uploadHandler))
	downloadZipHandlerWithCors := enableCors(http.HandlerFunc(downloadZipHandler))
	mux.Handle("/upload", uploadHandlerWithCors)
	mux.Handle("/download-all", downloadZipHandlerWithCors)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Println("Server started successfully.")
}
