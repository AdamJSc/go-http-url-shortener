package responseservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	headers    map[string]string
	statusCode int
}

func (r JSONResponse) Write(w http.ResponseWriter) JSONResponse {
	p := r.payload
	h := r.headers
	s := r.statusCode

	// set default http status
	if s == 0 {
		s = http.StatusInternalServerError
	}

	// set default content type header
	w.Header().Set("Content-Type", "application/json")

	// set header overrides
	for k, v := range h {
		w.Header().Set(k, v)
	}

	// set status code
	w.WriteHeader(s)

	// if payload is not blank, then write as body
	if p != (payload{}) {
		// parse response payload
		body, _ := json.Marshal(p)
		fmt.Fprintf(w, string(body))
	}

	return r
}

// NewEmptyResponse returns a new JSONResponse without a payload
func NewEmptyResponse(code int, headers ...string) JSONResponse {
	headersMap := map[string]string{}

	var k, v string
	i := 0

	// map our array of header strings in exact pairs
	// e.g. "Content-Type", "application/json", "Location", "http://bbc.co.uk"
	for (i + 1) < len(headers) {
		k = headers[i]
		v = headers[i+1]

		if k != "" && v != "" {
			headersMap[k] = v
		}

		i = i + 2
	}

	r := JSONResponse{
		statusCode: code,
		headers:    headersMap,
	}

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

// ParseJSON marshals a JSON payload into a map
func ParseJSON(resp *http.Response) map[string]interface{} {
	bytes, _ := ioutil.ReadAll(resp.Body)

	jsonMap := map[string]interface{}{}
	json.Unmarshal(bytes, &jsonMap)

	return jsonMap
}
