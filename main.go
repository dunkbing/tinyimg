package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/image"

	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

type RequestBody struct {
	Files []string `json:"files"`
}

func enableCors(next http.Handler) http.Handler {
	allowedOriginsMap := map[string]bool{}
	for _, v := range config.AllowedOrigins {
		allowedOriginsMap[v] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Request Headers: ", r.Header)
		origin := r.Header.Get("Origin")
		if allowedOriginsMap[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
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
	id := uuid.New()
	filename = fmt.Sprintf("%s_%s", id.String(), filename)

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
		http.Error(w, "File not found", http.StatusBadRequest)
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

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Change the map to hold values of the type visitor.
var visitors = make(map[string]*visitor)
var mu sync.Mutex

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	fmt.Println(ip)
	if !exists {
		limiter := rate.NewLimiter(30, 120)
		// Include the current time when creating a new visitor.
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}

func limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send hello message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/download-all", downloadZipHandler)
	mux.HandleFunc("/image", serveImgHandler)
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", enableCors((mux)))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
