package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"optipic/converter/config"
	"optipic/converter/image"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type RequestBody struct {
	Files []string `json:"files"`
}

var allowedOrigin = os.Getenv("ALLOWED_ORIGIN")

func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	startTime := time.Now()
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

	mimeType := http.DetectContentType(data)
	if !isImage(mimeType) {
		http.Error(w, "Invalid file format. Only images are allowed.", http.StatusBadRequest)
		return
	}
	fmt.Println("content type", mimeType, header.Filename)
	fileType, _ := image.GetFileType(mimeType)
	filename := filepath.Base(header.Filename)
	filename = strings.ReplaceAll(filename, " ", "_")

	ext := filepath.Ext(filename)
	filename = strings.Replace(filename, ext, fmt.Sprintf(".%s", fileType), 1)
	fmt.Println("ext", ext)
	ext = fmt.Sprintf(".%s", fileType)

	c := config.GetConfig()
	dest := filepath.Join(c.App.InDir, filename)
	slog.Info("Upload", "dest", dest)
	err = os.WriteFile(dest, data, 0644)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}
	took := time.Since(startTime).Seconds()
	fmt.Println("Write to file took", took)
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
	zippedPath, err := fm.ZipFiles(body.Files)
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

	// Set the content type header for zip files
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))

	// Copy the zip file content to the response writer
	_, err = io.Copy(w, zipFile)
	if err != nil {
		http.Error(w, "Error serving zip file", http.StatusInternalServerError)
		return
	}
}

func serveImgHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("f")
	if fileName == "" {
		http.Error(w, "Please provide a valid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join("output", fileName)

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Open and serve the file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Error opening file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the content type header based on the file extension
	contentType := getContentType(fileName)
	w.Header().Set("Content-Type", contentType)

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "Error serving file", http.StatusInternalServerError)
		return
	}
}

func getContentType(fileName string) string {
	switch filepath.Ext(fileName) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
}

func isImage(mimeType string) bool {
	mimeType = strings.ToLower(mimeType)
	return strings.HasPrefix(mimeType, "image/")
}

func main() {
	mux := http.NewServeMux()
	uploadHandlerWithCors := enableCors(http.HandlerFunc(uploadHandler))
	downloadZipHandlerWithCors := enableCors(http.HandlerFunc(downloadZipHandler))
	serveImageHandlerWithCors := enableCors(http.HandlerFunc(serveImgHandler))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send hello message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))
	mux.Handle("/upload", uploadHandlerWithCors)
	mux.Handle("/download-all", downloadZipHandlerWithCors)
	mux.Handle("/image", serveImageHandlerWithCors)
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
