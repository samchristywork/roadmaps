package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
)

func writeToFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	err = writeToFile("data/"+id+".json", body)
	if err != nil {
		http.Error(w, "Error writing to file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Data saved successfully")
}

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

	http.HandleFunc("/save", saveHandler)

	log.Println("Starting server at :1235")
	if err := http.ListenAndServe(":1235", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
