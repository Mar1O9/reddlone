package api

import (
	"encoding/json"
	"flag"
	"net/http"
	"sync"

	"github.com/Mar1O9/reddlone/views"
	"github.com/gorilla/mux"
)

var wg sync.WaitGroup

func routes() http.Handler {
	r := mux.NewRouter()

	var dir string

	flag.StringVar(&dir, "dir", "static", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.Handle("/", getIndex())
	return r
}
func getIndex() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		views.Home("hello").Render(r.Context(), w)
	})
}
