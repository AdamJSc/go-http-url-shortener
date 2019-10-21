package responseservice

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// payload represents response body
type payload struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

// JSONResponse represents status code and payload of a response
type JSONResponse struct {
	payload    payload
	statusCode int
}

func (r JSONResponse) Write(w http.ResponseWriter) JSONResponse {
	payload := r.payload
	statusCode := r.statusCode

	// set default http status
	if statusCode == 0 {
		statusCode = http.StatusInternalServerError
	}

	// parse response payload
	body, _ := json.Marshal(payload)

	// set headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// write body
	fmt.Fprintf(w, string(body))

	return r
}

// NewOkResponse returns a new JSONResponse representing a successful request
func NewOkResponse(data interface{}) JSONResponse {
	r := JSONResponse{
		payload: payload{
			Status: "ok",
			Data:   data,
		},
		statusCode: http.StatusOK,
	}

	return r
}

// NewErrResponse returns a new JSONResponse representing a request that generated an error
func NewErrResponse(message string, code ...int) JSONResponse {
	r := JSONResponse{
		payload: payload{
			Status: "err",
			Data: map[string]string{
				"message": message,
			},
		},
	}

	// set status code if provided to function call
	if len(code) > 0 {
		r.statusCode = code[0]
	}

	return r
}
