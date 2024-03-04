package main

import (
	"encoding/json"
	"github.com/dunkbing/tinyimg/converter/config"
	"github.com/dunkbing/tinyimg/converter/handlers"
	"log"
	"net/http"
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
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send hello message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Hello",
		})
	}))
	mux.HandleFunc("/upload", handlers.Upload)
	mux.HandleFunc("/download-all", handlers.DownloadAll)
	mux.HandleFunc("/image", handlers.ServeImg)
	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started on port 8080")
	err := http.ListenAndServe(":8080", enableCors(handlers.Limit(mux)))
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
