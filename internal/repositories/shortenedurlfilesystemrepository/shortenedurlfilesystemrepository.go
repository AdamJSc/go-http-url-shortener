package shortenedurlfilesystemrepository

import (
	"encoding/json"
	"errors"
	"fmt"
	"http-url-shortener/internal/entities/shortenedurl"
	"io/ioutil"
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

// Create a new Shortened URL on file system
func (f FileSystem) Create(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	if u.GetLong() == "" || u.GetShort() == "" {
		// nothing to save
		return shortenedurl.ShortenedURL{}, errors.New("Shortened URL is empty")
	}

	path := getPathToDbFile(f)
	m := loadManifest(path)

	if m[u.GetLong()] != "" {
		// already exists
		return shortenedurl.ShortenedURL{}, errors.New("Shortened URL already exists")
	}

	m[u.GetLong()] = u.GetShort()
	if saveManifest(path, m) == false {
		// unable to save
		return shortenedurl.ShortenedURL{}, errors.New("Shortened URL could not be created")
	}

	return u, nil
}

// RetrieveByShortCode retrieves a Shortened URL by its short code
func (f FileSystem) RetrieveByShortCode(shortcode string) (shortenedurl.ShortenedURL, error) {
	m := loadManifest(getPathToDbFile(f))

	// try to retrieve by URL's short code
	for l, s := range m {
		if s == shortcode {
			return shortenedurl.New(l, s), nil
		}
	}

	// no matching manifest entries
	return shortenedurl.ShortenedURL{}, errors.New("Shortened URL does not exist")
}

// RetrieveByLongURL retrieves a Shortened URL by its origin (long) URL
func (f FileSystem) RetrieveByLongURL(longURL string) (shortenedurl.ShortenedURL, error) {
	m := loadManifest(getPathToDbFile(f))

	// try to retrieve by origin (long) URL
	if m[longURL] != "" {
		return shortenedurl.New(longURL, m[longURL]), nil
	}

	// no matching manifest entries
	return shortenedurl.ShortenedURL{}, errors.New("Shortened URL does not exist")
}

func getPathToDbFile(f FileSystem) string {
	return f.basePath + "/db.txt"
}

func loadManifest(path string) map[string]string {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		return map[string]string{}
	}

	m := map[string]string{}

	err = json.Unmarshal(fileContents, &m)
	if err != nil {
		fmt.Println(err)
		return map[string]string{}
	}

	return m
}

func saveManifest(path string, data map[string]string) bool {
	fileContents, err := json.Marshal(data)
	if err != nil {
		return false
	}

	err = ioutil.WriteFile(path, fileContents, 0644)
	if err != nil {
		return false
	}

	return true
}
