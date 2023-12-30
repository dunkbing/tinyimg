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

const allowedOrigin = "http://localhost:8000, https://tinyimg.deno.dev"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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

	formatStr := r.FormValue("formats")
	formats := strings.Split(formatStr, ",")

	filename := header.Filename
	ext := filepath.Ext(header.Filename)
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
	fileManager.HandleFile(&f)
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
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download-all", downloadZipHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Println("Server started successfully.")
}
