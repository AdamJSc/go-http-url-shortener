package main

import (
	"http-url-shortener/internal/handlers"
	"http-url-shortener/internal/repositories/shortenedurlfilesystemrepository"
	"log"
	"net/http"
	"os"
)

func main() {
	workdir, _ := os.Getwd()
	repository := shortenedurlfilesystemrepository.New(workdir + "/data")

	http.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handlers.PostShorten(repository, w, r).Write(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
