package shortenedurlfilesystemrepository

import (
	"errors"
	"http-url-shortener/internal/entities/shortenedurl"
	"os"
)

// FileSystem represents a file system to perform operations on
type FileSystem struct {
	basePath string
}

// New instance of FileSystem type
func New(p string) FileSystem {
	return FileSystem{
		basePath: p,
	}
}

// Create a new URL on file system
func (f FileSystem) Create(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	fullPath := getPathToURLFile(f, u)
	_, err := os.OpenFile(fullPath, os.O_RDONLY, 0644)

	if err == nil {
		// file already exists...
		return shortenedurl.ShortenedURL{}, errors.New("Shortened URL already exists")
	}

	return u, nil
}

// Retrieve an existing URL from file system
func (f FileSystem) Retrieve(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	// @TODO implement
	return u, nil
}

// Update an existing URL from file system
func (f FileSystem) Update(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	// @TODO implement
	return u, nil
}

// Delete an existing URL from file system
func (f FileSystem) Delete(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	// @TODO implement
	return u, nil
}

func getPathToURLFile(f FileSystem, u shortenedurl.ShortenedURL) string {
	return f.basePath + "/" + u.GetShort() + ".txt"
}
