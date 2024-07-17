package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dunkbing/tinyimg/tinyimg/config"
	"github.com/dunkbing/tinyimg/tinyimg/handlers"
)

func enableCors(next http.Handler) http.Handler {
	allowedOriginsMap := map[string]bool{}
	for _, v := range config.AllowedOrigins {
		allowedOriginsMap[v] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowedOriginsMap[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	handler := handlers.New()
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "pong")
	}))
	mux.HandleFunc("POST /upload", handler.Upload)
	mux.HandleFunc("POST /download-all", handler.DownloadAll)
	mux.HandleFunc("/image", handler.ServeImg)
	mux.HandleFunc("/video", handler.ServeVideo)
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", enableCors(handlers.Limit(mux)))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
