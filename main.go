package main

import (
	"fmt"
	"io"
	"net/http"
	"optipic/converter/image"
	"path/filepath"
)

type RequestBody struct {
    Keys []string `json:"keys"`
}

var bucketName = "optipic"
var accountId = "<account_id>"
var accessKeyId = "<access_key_id>"
var accessKeySecret = "<access_key_secret>"

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

    id := r.FormValue("id")
    f := image.File{
        Data:     data,
        Ext:      filepath.Ext(header.Filename),
        ID:       id,
        MimeType: header.Header.Get("Content-Type"),
        Name:     header.Filename,
        Size:     header.Size,
    }

    fmt.Fprintf(w, "File uploaded successfully: %+v", f)
}

func main() {
    http.HandleFunc("/upload", uploadHandler)

    http.ListenAndServe(":8080", nil)
}
