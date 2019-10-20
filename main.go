package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// responsePayload represents a standardised JSON payload returned by API
type responsePayload struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		payload := responsePayload{
			Status: "ok",
			Data: map[string]string{
				"hello": "world",
			},
		}

		writeResponse(w, payload, http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeResponse(w http.ResponseWriter, p responsePayload, code int) {
	// parse response payload
	body, _ := json.Marshal(p)

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// write body
	fmt.Fprintf(w, string(body))
}
