package main

import (
	"http-url-shortener/internal/services/responseservice"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"hello": "world",
		}

		responseservice.NewOkResponse(data).Write(w)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
