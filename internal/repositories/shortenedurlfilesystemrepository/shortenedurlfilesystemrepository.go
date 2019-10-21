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

// Retrieve an existing Shortened URL from file system
func (f FileSystem) Retrieve(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error) {
	m := loadManifest(getPathToDbFile(f))

	// try to retrieve by URL's long address
	if u.GetLong() != "" && m[u.GetLong()] != "" {
		return shortenedurl.New(u.GetLong(), m[u.GetLong()]), nil
	}

	if u.GetShort() != "" {
		// try to retrieve by URL's short address
		for long, short := range m {
			if u.GetShort() == short {
				return shortenedurl.New(long, short), nil
			}
		}
	}

	// no matching manifest entries
	return shortenedurl.ShortenedURL{}, errors.New("Shortened URL does not exist")
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
