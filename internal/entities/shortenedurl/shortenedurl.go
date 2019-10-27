package shortenedurl

// ShortenedURL type represents a URL to be handled by the system
type ShortenedURL struct {
	long  string
	short string
}

// New creates a new instance of type ShortenedURL
func New(l string, s string) ShortenedURL {
	return ShortenedURL{
		long:  l,
		short: s,
	}
}

// GetLong retrieves value of ShortenedURL instance's `long` property
func (u ShortenedURL) GetLong() string {
	return u.long
}

// GetShort retrieves value of ShortenedURL instance's `short` property
func (u ShortenedURL) GetShort() string {
	return u.short
}
