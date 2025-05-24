package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func writeToFile(filename string, data []byte) error {
	if err := os.MkdirAll("data", 0755); err != nil {
		return err
	}
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" || strings.ContainsAny(id, "/\\") || strings.Contains(id, "..") {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := os.Remove("data/" + id + ".json"); err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error deleting file", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Deleted")
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" || strings.ContainsAny(id, "/\\") || strings.Contains(id, "..") {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	if !json.Valid(body) {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
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

func treesHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := os.ReadDir("data")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintln(w, `<!DOCTYPE html><html><head><meta charset="UTF-8">
<title>Roadmaps</title>
<style>
  body { font-family: monospace; max-width: 600px; margin: 2rem auto; padding: 0 1rem; }
  h1 { font-size: 1.2rem; }
  ul { list-style: none; padding: 0; }
  li { margin: 0.4rem 0; display: flex; align-items: center; gap: 0.5rem; }
  a { color: #007bff; text-decoration: none; }
  a:hover { text-decoration: underline; }
  .del { background: none; border: none; color: #aaa; cursor: pointer; font-size: 14px; padding: 0 4px; }
  .del:hover { color: #dc3545; }
  .new { margin-top: 1.5rem; display: flex; gap: 0.5rem; }
  .new input { font-family: monospace; padding: 4px 8px; border: 1px solid #ccc; border-radius: 4px; }
  .new button { padding: 4px 12px; cursor: pointer; }
</style></head><body>
<h1>Roadmaps</h1><ul>`)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".json") {
				id := strings.TrimSuffix(e.Name(), ".json")
				fmt.Fprintf(w, `<li><a href="/?id=%s">%s</a><button class="del" data-id="%s" onclick="del(this.dataset.id)" title="Delete">✕</button></li>`+"\n",
					url.QueryEscape(id), html.EscapeString(id), html.EscapeString(id))
			}
		}
	}
	fmt.Fprintln(w, `</ul>
<div class="new">
  <input id="newId" type="text" placeholder="new roadmap name">
  <button onclick="go()">Create</button>
</div>
<script>
  document.getElementById('newId').addEventListener('keydown', e => { if (e.key==='Enter') go(); });
  function go() {
    const id = document.getElementById('newId').value.trim();
    if (id) window.location.href = '/?id=' + encodeURIComponent(id);
  }
  function del(id) {
    if (!confirm('Delete "' + id + '"?')) return;
    fetch('/delete?id=' + encodeURIComponent(id), { method: 'DELETE' })
      .then(r => { if (r.ok) location.reload(); else alert('Delete failed'); });
  }
</script></body></html>`)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested URL: %s from %s", r.URL.Path, r.RemoteAddr)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("id") == "" {
			http.Redirect(w, r, "/trees", http.StatusFound)
			return
		}
		log.Printf("Serving index.html to %s", r.RemoteAddr)
		http.ServeFile(w, r, "index.html")
	})

	fs := http.FileServer(http.Dir("./data/"))
	http.Handle("/data/", http.StripPrefix("/data", logRequest(fs)))

	http.HandleFunc("/trees", treesHandler)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/delete", deleteHandler)

	log.Println("Starting server at :1235")
	if err := http.ListenAndServe(":1235", nil); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
