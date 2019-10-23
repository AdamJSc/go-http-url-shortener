package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

type apiResponse struct {
	Status string            `json:"status"`
	Data   map[string]string `json:"data"`
}

func TestItReceivesANotFoundStatusWhenRootIsRequested(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusNotFound, resp.StatusCode))
	}
}

func TestItReceivesANotFoundStatusWhenNonExistentShortCodeIsRequested(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/ABCD", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusNotFound, resp.StatusCode))
	}
}

func TestItReceivesAMethodNotAllowedStatusWhenInvalidMethodIsUsed(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusMethodNotAllowed, resp.StatusCode))
	}
}

func TestItFailsToShortenAURLWhenPayloadIsMissing(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/api/shorten", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusBadRequest, resp.StatusCode))
	}

	json := parseJSON(resp)

	if json.Status != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json.Status))
	}

	if json.Data["message"] != "unexpected end of JSON input" {
		t.Error(fmt.Sprintf("Expected message of 'unexpected end of JSON input', instead received '%s'", json.Data["message"]))
	}
}

func TestItFailsToShortenAURLWhenURLIsMissingFromPayload(t *testing.T) {
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

	json := parseJSON(resp)

	if json.Status != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json.Status))
	}

	if json.Data["message"] != "`url` is a non-string or missing" {
		t.Error(fmt.Sprintf("Expected message of '`url` is a non-string or missing', instead received '%s'", json.Data["message"]))
	}
}

func TestItFailsToShortenAURLWhenURLIsNotAValidType(t *testing.T) {
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

	json := parseJSON(resp)

	if json.Status != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json.Status))
	}

	if json.Data["message"] != "`url` is a non-string or missing" {
		t.Error(fmt.Sprintf("Expected message of '`url` is a non-string or missing', instead received '%s'", json.Data["message"]))
	}
}

func TestItFailsToShortenAURLWhenURLIsNotAValidFormat(t *testing.T) {
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

	json := parseJSON(resp)

	if json.Status != "err" {
		t.Error(fmt.Sprintf("Expected JSON status of 'err', instead received '%s'", json.Status))
	}

	if json.Data["message"] != "`url` is not a valid URL" {
		t.Error(fmt.Sprintf("Expected message of '`url` is not a valid URL', instead received '%s'", json.Data["message"]))
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

	json := parseJSON(resp)

	if json.Status != "ok" {
		t.Error(fmt.Sprintf("Expected JSON status of 'ok', instead received '%s'", json.Status))
	}

	expectedShortURLPrefix := "http://localhost:8080/"
	if strings.HasPrefix(json.Data["shortURL"], expectedShortURLPrefix) != true {
		t.Error(fmt.Sprintf(
			"Expected shortURL to begin with of '%s', instead shortURL is '%s'",
			expectedShortURLPrefix,
			json.Data["shortURL"],
		))
	}

	shortCode := strings.TrimPrefix(json.Data["shortURL"], expectedShortURLPrefix)
	if len(shortCode) != 4 {
		t.Error(fmt.Sprintf(
			"Expected shortURL code to be 4 characters long, instead code length is %d ('%s')",
			len(shortCode),
			json.Data["shortURL"],
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

	json := parseJSON(resp)

	if json.Status != "ok" {
		t.Error(fmt.Sprintf("Expected JSON status of 'ok', instead received '%s'", json.Status))
	}

	expectedShortURL := "http://localhost:8080/ABC1"
	if json.Data["shortURL"] != expectedShortURL {
		t.Error(fmt.Sprintf("Expected shortURL of '%s', instead received '%s'", expectedShortURL, json.Data["shortURL"]))
	}

	// clean up
	clearTestData()
}

func parseJSON(resp *http.Response) apiResponse {
	bytes, _ := ioutil.ReadAll(resp.Body)

	jsonMap := apiResponse{}
	json.Unmarshal(bytes, &jsonMap)

	return jsonMap
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
