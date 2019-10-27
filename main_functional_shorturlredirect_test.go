package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestItReturnsNotFoundWhenRootIsRequested(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusNotFound, resp.StatusCode))
	}
}

func TestItReturnsMethodNotAllowedWhenInvalidMethodIsUsed(t *testing.T) {
	r := httptest.NewRequest("POST", "http://localhost:8080/", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusMethodNotAllowed, resp.StatusCode))
	}
}

func TestItReturnsNotFoundWhenNoPreviousShortenedURLExists(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost:8080/ABCD", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusNotFound, resp.StatusCode))
	}
}

func TestItReturnsNotFoundWhenURLShortCodeDoesNotExist(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	r := httptest.NewRequest("GET", "http://localhost:8080/DEF2", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusNotFound, resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if len(string(body)) != 0 {
		t.Error(fmt.Sprintf("Expected empty body, instead received '%s'", body))
	}

	// clean up
	clearTestData()
}

func TestItReturnsRedirectWhenURLShortCodeExists(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	r := httptest.NewRequest("GET", "http://localhost:8080/ABC1", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusMovedPermanently {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusMovedPermanently, resp.StatusCode))
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if len(string(body)) != 0 {
		t.Error(fmt.Sprintf("Expected empty body, instead received '%s'", body))
	}

	location := resp.Header.Get("Location")
	if location != "http://bbc.co.uk" {
		t.Error(fmt.Sprintf(
			"Expected location header '%s', instead received '%s'",
			"http://bbc.co.uk",
			location,
		))
	}

	// clean up
	clearTestData()
}

func TestItReturnsMethodNotAllowedWhenInvalidMethodIsUsedOnExistingShortCode(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	r := httptest.NewRequest("POST", "http://localhost:8080/ABC1", nil)
	w := httptest.NewRecorder()

	apiHandler(w, r)
	resp := w.Result()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Sprintf("Expected status code %d, instead received %d", http.StatusMethodNotAllowed, resp.StatusCode))
	}

	// clean up
	clearTestData()
}
