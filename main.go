package main

import (
	"http-url-shortener/internal/handlers"
	"http-url-shortener/internal/repositories/shortenedurlfilesystemrepository"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", apiHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	workdir, _ := os.Getwd()
	repository := shortenedurlfilesystemrepository.New(workdir + "/data")

	// shortener endpoint
	if r.URL.Path == "/api/shorten" {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		handlers.PostShorten(repository, w, r).Write(w)
		return
	}

	// try to redirect a short URL
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handlers.GetShortURLRedirect(repository, w, r).Write(w)
}
