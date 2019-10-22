package handlers

import (
	"encoding/json"
	"errors"
	"http-url-shortener/internal/entities/shortenedurl"
	"http-url-shortener/internal/repositories/shortenedurlfilesystemrepository"
	"http-url-shortener/internal/services/responseservice"
	"http-url-shortener/internal/services/shortcodeservice"
	"io/ioutil"
	"net/http"
	"net/url"
)

// PostShorten handles request to shorten a URL
func PostShorten(
	fs shortenedurlfilesystemrepository.FileSystem,
	w http.ResponseWriter,
	r *http.Request,
) responseservice.JSONResponse {
	// extract URL property from request body
	urlValue, err := getValueOfURLFromRequestBody(r)
	if err != nil {
		return responseservice.NewErrResponse(err.Error(), http.StatusBadRequest)
	}

	// check if we've already shortened it
	existingURL, err := fs.Retrieve(shortenedurl.New(urlValue, ""))
	if err == nil {
		// return our existing record
		return responseservice.NewOkResponse(map[string]string{
			"long":  existingURL.GetLong(),
			"short": existingURL.GetShort(),
		})
	}

	// URL is new, let's generate a new shortcode
	shortCode := shortcodeservice.Generate()

	// loop until we have a unique short code...
	for err != nil {
		shortCode = shortcodeservice.Generate()
		_, err = fs.Retrieve(shortenedurl.New("", shortCode))
	}

	// save our shortened URL
	shortenedURL, err := fs.Create(shortenedurl.New(urlValue, shortCode))
	if err != nil {
		return responseservice.NewErrResponse(err.Error())
	}

	// return our new record
	return responseservice.NewOkResponse(map[string]string{
		"long":  shortenedURL.GetLong(),
		"short": shortenedURL.GetShort(),
	})
}

// GetShortURLRedirect handles request to redirect a short URL
func GetShortURLRedirect(
	fs shortenedurlfilesystemrepository.FileSystem,
	w http.ResponseWriter,
	r *http.Request,
) {
	// default to 404
	w.WriteHeader(http.StatusNotFound)
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
