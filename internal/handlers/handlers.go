package handlers

import (
	"http-url-shortener/internal/repositories/shortenedurlfilesystemrepository"
	"http-url-shortener/internal/services/responseservice"
	"net/http"
)

// PostShorten handles request to shorten a URL
func PostShorten(fs shortenedurlfilesystemrepository.FileSystem, w http.ResponseWriter, r *http.Request) {
	responseservice.NewOkResponse(map[string]string{"hello": "world"}).Write(w)
}
