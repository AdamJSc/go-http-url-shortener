package main

import (
	"encoding/json"
	"fmt"
	"http-url-shortener/internal/services/responseservice"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"hello": "world",
		}

		writeOkResponse(w, data)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// writeOk writes a standardised JSON response for a successful request
func writeOkResponse(w http.ResponseWriter, data interface{}) {
	payload := responseservice.NewPayload("ok", data)

	writeResponse(w, payload, http.StatusOK)
}

// writeErr writes a standardised JSON response for a failed request
func writeErrResponse(w http.ResponseWriter, data interface{}, code ...int) {
	payload := responseservice.NewPayload("err", data)

	// default to 500 unless an alternative code has been supplied
	statusCode := http.StatusInternalServerError
	if len(code) != 0 {
		statusCode = code[0]
	}

	writeResponse(w, payload, statusCode)
}

func writeResponse(w http.ResponseWriter, p responseservice.Payload, code int) {
	// parse response payload
	body, _ := json.Marshal(p)

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	// write body
	fmt.Fprintf(w, string(body))
}
