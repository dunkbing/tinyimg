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
	Keys []string `json:"keys"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

	id := r.FormValue("id")
	filename := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	ext := filepath.Ext(header.Filename)
	f := image.File{
		Data:     data,
		Ext:      ext,
		ID:       id,
		MimeType: header.Header.Get("Content-Type"),
		Name:     filename,
		Size:     header.Size,
	}

	if !isImage(f.MimeType) {
		http.Error(w, "Invalid file format. Only images are allowed.", http.StatusBadRequest)
		return
	}

	fileManager := image.NewFileManager()
	fileManager.HandleFile(&f)
	stat, errs := fileManager.Convert()
	if len(errs) > 0 {
		http.Error(w, "Failed to convert file", http.StatusInternalServerError)
		return
	}

	// Success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stat)
}

func isImage(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	return strings.HasPrefix(mimeType, "image/")
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
	log.Println("Server started successfully.")
}
