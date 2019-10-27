package repositoryinterface

import "http-url-shortener/internal/entities/shortenedurl"

// RepositoryInterface defines interface for a Shortened URL repository
type RepositoryInterface interface {
	Create(u shortenedurl.ShortenedURL) (shortenedurl.ShortenedURL, error)
	RetrieveByShortCode(shortcode string) (shortenedurl.ShortenedURL, error)
	RetrieveByLongURL(longURL string) (shortenedurl.ShortenedURL, error)
}
