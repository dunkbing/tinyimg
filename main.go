package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/handlers"
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
	mux.Handle("/ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprint(w, "pong")
	}))
	mux.HandleFunc("POST /upload", handlers.Upload)
	mux.HandleFunc("POST /download-all", handlers.DownloadAll)
	mux.HandleFunc("/image", handlers.ServeImg)
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", enableCors(handlers.Limit(mux)))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
