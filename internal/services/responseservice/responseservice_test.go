package responseservice

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItReturnsANewEmptyResponseWithNoHeaders(t *testing.T) {
	response := NewEmptyResponse(http.StatusMovedPermanently)

	if response.statusCode != http.StatusMovedPermanently {
		t.Errorf("Expected status code of %d, instead received %d", http.StatusMovedPermanently, response.statusCode)
	}

	if response.payload != (payload{}) {
		t.Errorf("Expected blank payload, instead received '%+v'", response.payload)
	}

	if len(response.headers) != 0 {
		t.Errorf("Expected blank headers, instead received '%+v'", response.headers)
	}
}

func TestItReturnsANewEmptyResponseWithSomeHeaders(t *testing.T) {
	response := NewEmptyResponse(http.StatusMovedPermanently, "Location", "http://bbc.co.uk")

	if response.statusCode != http.StatusMovedPermanently {
		t.Errorf("Expected status code of %d, instead received %d", http.StatusMovedPermanently, response.statusCode)
	}

	if response.payload != (payload{}) {
		t.Errorf("Expected blank payload, instead received '%+v'", response.payload)
	}

	if len(response.headers) != 1 {
		t.Errorf("Expected 1 header, instead received '%+v'", response.headers)
	}

	if response.headers["Location"] != "http://bbc.co.uk" {
		t.Errorf("Expected location header of '%s', instead received '%+s'", "http://bbc.co.uk", response.headers["Location"])
	}
}

func TestItReturnsANewOkResponse(t *testing.T) {
	response := NewOkResponse(map[string]string{"hello": "world"})

	if response.statusCode != http.StatusOK {
		t.Errorf("Expected status code of %d, instead received %d", http.StatusOK, response.statusCode)
	}

	if response.payload.Status != "ok" {
		t.Errorf("Expected payload status of '%s', instead received '%s'", "ok", response.payload.Status)
	}

	if response.payload.Data.(map[string]string)["hello"] != "world" {
		t.Errorf("Expected payload data of {\"hello\":\"world\"}, instead received '%+v'", response.payload.Data)
	}

	if len(response.headers) != 0 {
		t.Errorf("Expected blank headers, instead received '%+v'", response.headers)
	}
}

func TestItReturnsANewErrResponse(t *testing.T) {
	response := NewErrResponse("Feels badgateway man :(", http.StatusBadGateway)

	if response.statusCode != http.StatusBadGateway {
		t.Errorf("Expected status code of %d, instead received %d", http.StatusBadGateway, response.statusCode)
	}

	if response.payload.Status != "err" {
		t.Errorf("Expected payload status of '%s', instead received '%s'", "err", response.payload.Status)
	}

	if response.payload.Data.(map[string]string)["message"] != "Feels badgateway man :(" {
		t.Errorf("Expected payload data of 'Feels badgateway man :(', instead received '%+v'", response.payload.Data)
	}

	if len(response.headers) != 0 {
		t.Errorf("Expected blank headers, instead received '%+v'", response.headers)
	}
}

func TestItSuccessfullyWritesAResponse(t *testing.T) {
	response := JSONResponse{
		payload: payload{
			Status: "ok",
			Data:   map[string]string{"hello": "world"},
		},
		headers:    map[string]string{"X-My-Test-Header-1": "ABC123", "X-My-Test-Header-2": "DEF456"},
		statusCode: http.StatusTeapot,
	}

	writer := httptest.NewRecorder()
	response.Write(writer)

	// check for expected status code
	if writer.Code != http.StatusTeapot {
		t.Errorf("Expected status code of %d, instead received %d", http.StatusTeapot, writer.Code)
	}

	// check for expected headers
	if writer.Header().Get("X-My-Test-Header-1") != "ABC123" {
		t.Errorf(
			"Expected value of '%s' header to be '%s', instead received '%s'",
			"X-My-Test-Header-1",
			"ABC123",
			writer.Header().Get("X-My-Test-Header-1"),
		)
	}

	if writer.Header().Get("X-My-Test-Header-2") != "DEF456" {
		t.Errorf(
			"Expected value of '%s' header to be '%s', instead received '%s'",
			"X-My-Test-Header-2",
			"DEF456",
			writer.Header().Get("X-My-Test-Header-2"),
		)
	}

	if writer.Header().Get("Content-Type") != "application/json" {
		t.Errorf(
			"Expected value of '%s' header to be '%s', instead received '%s'",
			"Content-Type",
			"application/json",
			writer.Header().Get("Content-Type"),
		)
	}

	if len(writer.Header()) != 3 {
		t.Errorf(
			"Expected 3 headers, instead received '%+v'",
			writer.Header(),
		)
	}

	// check for expected payload
	payload := ParseJSON(writer.Result())

	if payload["status"] != "ok" {
		t.Errorf("Expected payload status of '%s', instead received '%s'", "ok", payload["status"])
	}

	if payload["data"].(map[string]interface{})["hello"].(string) != "world" {
		t.Errorf("Expected data payload of {\"hello\":\"world\"}, instead received '%+v'", payload["data"])
	}
}
