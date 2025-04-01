package main

import (
	"log"
	"net/http"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested URL: %s from %s", r.URL.Path, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving index.html to %s", r.RemoteAddr)
		http.ServeFile(w, r, "index.html")
	})

	fs := http.FileServer(http.Dir("./data/"))
	http.Handle("/data/", http.StripPrefix("/data", logRequest(fs)))

	log.Println("Starting server at :1235")
	if err := http.ListenAndServe(":1235", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
