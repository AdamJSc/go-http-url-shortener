package url

// URL type represents a URL to be handled by the system
type URL struct {
	long  string
	short string
}

// New creates a new instance of type URL
func New(l string, s string) URL {
	return URL{
		long:  l,
		short: s,
	}
}

// GetLong retrieves value of URL instance's `long` property
func (u URL) GetLong() string {
	return u.long
}

// GetShort retrieves value of URL instance's `short` property
func (u URL) GetShort() string {
	return u.short
}
