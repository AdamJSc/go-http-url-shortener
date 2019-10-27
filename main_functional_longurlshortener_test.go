package main

import (
	"fmt"
	"http-url-shortener/internal/services/responseservice"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestItFailsToShortenAURLWhenPayloadIsEmpty(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/shorten", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusBadRequest, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	if jsonData["message"] != "unexpected end of JSON input" {
		t.Error(fmt.Sprintf("Expected message of 'unexpected end of JSON input', instead received '%s'", jsonData["message"]))
	}
}

func TestItFailsToShortenAURLWhenLongURLIsMissingFromPayload(t *testing.T) {
	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"incorrect_url_key": "http://bbc.co.uk"}`),
	)
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusBadRequest, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	if jsonData["message"] != "`url` is a non-string or missing" {
		t.Error(fmt.Sprintf("Expected message of '`url` is a non-string or missing', instead received '%s'", jsonData["message"]))
	}
}

func TestItFailsToShortenAURLWhenLongURLIsNotAValidType(t *testing.T) {
	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"url": 123}`),
	)
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusBadRequest, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	if jsonData["message"] != "`url` is a non-string or missing" {
		t.Error(fmt.Sprintf("Expected message of '`url` is a non-string or missing', instead received '%s'", jsonData["message"]))
	}
}

func TestItFailsToShortenAURLWhenLongURLIsNotAValidFormat(t *testing.T) {
	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"url": "http//bbc.co.uk"}`),
	)
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusBadRequest, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	if jsonData["message"] != "`url` is not a valid URL" {
		t.Error(fmt.Sprintf("Expected message of '`url` is not a valid URL', instead received '%s'", jsonData["message"]))
	}
}

func TestItSuccessfullyReturnsANewShortURLWhenLongURLHasNotAlreadyBeenShortened(t *testing.T) {
	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"url": "http://bbc.co.uk"}`),
	)
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusOK, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "ok" {
		t.Error(fmt.Sprintf("Expected JSON status of 'ok', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	expectedShortURLPrefix := "http://localhost:8080/"
	if strings.HasPrefix(jsonData["shortURL"].(string), expectedShortURLPrefix) != true {
		t.Error(fmt.Sprintf(
			"Expected shortURL to begin with of '%s', instead shortURL is '%s'",
			expectedShortURLPrefix,
			jsonData["shortURL"],
		))
	}

	shortCode := strings.TrimPrefix(jsonData["shortURL"].(string), expectedShortURLPrefix)
	if len(shortCode) != 4 {
		t.Error(fmt.Sprintf(
			"Expected shortURL code to be 4 characters long, instead code length is %d ('%s')",
			len(shortCode),
			jsonData["shortURL"],
		))
	}

	// clean up
	clearTestData()
}

func TestItReturnsAKnownShortURLWhenLongURLHasAlreadyBeenShortened(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk":"ABC1"}`)

	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/api/shorten",
		strings.NewReader(`{"url": "http://bbc.co.uk"}`),
	)
	r.Header = map[string][]string{
		"Content-Type": {"application/json"},
	}

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusOK, resp.StatusCode))
	}

	json := responseservice.ParseJSON(resp)

	if json["status"] != "ok" {
		t.Error(fmt.Sprintf("Expected JSON status of 'ok', instead received '%s'", json["status"]))
	}

	jsonData := json["data"].(map[string]interface{})
	expectedShortURL := "http://localhost:8080/ABC1"
	if jsonData["shortURL"] != expectedShortURL {
		t.Error(fmt.Sprintf("Expected shortURL of '%s', instead received '%s'", expectedShortURL, jsonData["shortURL"]))
	}

	// clean up
	clearTestData()
}

func setTestData(data string) {
	clearTestData()
	ioutil.WriteFile(getTestDataPath(), []byte(data), 0644)
}

func clearTestData() {
	os.Remove(getTestDataPath())
}

func getTestDataPath() string {
	dir, _ := os.Getwd()
	return dir + "/data/db.txt"
}
