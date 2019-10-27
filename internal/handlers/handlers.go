package handlers

import (
	"encoding/json"
	"errors"
	"http-url-shortener/internal/entities/shortenedurl"
	"http-url-shortener/internal/repositories/repositoryinterface"
	"http-url-shortener/internal/services/responseservice"
	"http-url-shortener/internal/services/shortcodeservice"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// PostShorten handles request to shorten a URL
func PostShorten(
	repo repositoryinterface.RepositoryInterface,
	w http.ResponseWriter,
	r *http.Request,
) responseservice.JSONResponse {
	// extract URL property from request body
	urlValue, err := getValueOfURLFromRequestBody(r)
	if err != nil {
		return responseservice.NewErrResponse(err.Error(), http.StatusBadRequest)
	}

	// check if we've already shortened it
	existing, err := repo.RetrieveByLongURL(urlValue)
	if err == nil {
		// return our existing record
		return responseservice.NewOkResponse(map[string]string{
			"shortURL": "http://" + r.Host + "/" + existing.GetShort(),
		})
	}

	// URL is new, let's generate a new shortcode
	shortCode := shortcodeservice.Generate()

	// loop until we have a unique short code...
	for err == nil {
		shortCode = shortcodeservice.Generate()
		_, err = repo.RetrieveByShortCode(shortCode)
	}

	// save our shortened URL
	shortened, err := repo.Create(shortenedurl.New(urlValue, shortCode))
	if err != nil {
		return responseservice.NewErrResponse(err.Error())
	}

	// return our new record
	return responseservice.NewOkResponse(map[string]string{
		"shortURL": "http://" + r.Host + "/" + shortened.GetShort(),
	})
}

// GetShortURLRedirect handles request to redirect a short URL
func GetShortURLRedirect(
	repo repositoryinterface.RepositoryInterface,
	w http.ResponseWriter,
	r *http.Request,
) responseservice.JSONResponse {
	pathParts := strings.SplitN(r.URL.Path, "/", 2)
	if len(pathParts) < 2 {
		// root path "/" (no short code supplied)
		return responseservice.NewEmptyResponse(http.StatusNotFound)
	}

	shortCode := pathParts[1]
	shortenedURL, err := repo.RetrieveByShortCode(shortCode)
	if err == nil {
		// set redirect header to short code's corresponding long URL
		return responseservice.NewEmptyResponse(
			http.StatusMovedPermanently,
			"Location",
			shortenedURL.GetLong(),
		)
	}

	// nothing found
	return responseservice.NewEmptyResponse(http.StatusNotFound)
}

func getValueOfURLFromRequestBody(r *http.Request) (string, error) {
	// read request body
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	// parse request body as json
	var jsonBody map[string]interface{}
	err = json.Unmarshal(requestBody, &jsonBody)
	if err != nil {
		return "", err
	}

	// check that url exists in payload and is a string
	switch jsonBody["url"].(type) {
	case string:
		// ok
	default:
		return "", errors.New("`url` is a non-string or missing")
	}

	// check that URL is valid
	urlValue := jsonBody["url"].(string)
	_, err = url.ParseRequestURI(urlValue)
	if err != nil {
		return "", errors.New("`url` is not a valid URL")
	}

	return urlValue, nil
}
