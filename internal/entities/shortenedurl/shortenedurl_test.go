package shortenedurl

import (
	"reflect"
	"testing"
)

func TestItSuccessfullyReturnsAShortenedURL(t *testing.T) {
	long := "http://bbc.co.uk"
	short := "ABC1"

	result := New(long, short)
	resultType := reflect.TypeOf(result).Name()

	if resultType != "ShortenedURL" {
		t.Errorf("Expected object of type '%s', instead received '%s'", "ShortenedURL", resultType)
	}

	if result.GetLong() != long {
		t.Errorf("Expected long value of '%s', instead received '%s'", long, result.GetLong())
	}

	if result.GetShort() != short {
		t.Errorf("Expected short value of '%s', instead received '%s'", short, result.GetShort())
	}
}
