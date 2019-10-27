package shortenedurlfilesystemrepository

import (
	"http-url-shortener/internal/entities/shortenedurl"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestItCreatesANewFileSystemRepository(t *testing.T) {
	fs := New("/my/path")

	if reflect.TypeOf(fs).Name() != "FileSystem" {
		t.Errorf("Expected type of '%s', instead received '%s'", "FileSystem", reflect.TypeOf(fs).Name())
	}

	if fs.basePath != "/my/path" {
		t.Errorf("Expected basePath value of '%s', instead received '%s'", "/my/path", fs.basePath)
	}
}

func TestItFailsToCreateAShortenedURLIfSuppliedObjectHasNoValues(t *testing.T) {
	fs := getTestFsRepository()

	noShortCode := shortenedurl.New("longURL", "")
	expectedErrorMessage := "Shortened URL is empty"

	_, err := fs.Create(noShortCode)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message of '%s', instead received '%s'", expectedErrorMessage, err.Error())
	}

	noLongURL := shortenedurl.New("", "shortCode")
	expectedErrorMessage = "Shortened URL is empty"

	_, err = fs.Create(noLongURL)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message of '%s', instead received '%s'", expectedErrorMessage, err.Error())
	}

	noShortCodeOrLongURL := shortenedurl.New("", "")
	expectedErrorMessage = "Shortened URL is empty"

	_, err = fs.Create(noShortCodeOrLongURL)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message of '%s', instead received '%s'", expectedErrorMessage, err.Error())
	}
}

func TestItFailsToCreateAShortenedURLIfLongURLAlreadyExists(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	longURLAlreadyShortened := shortenedurl.New("http://bbc.co.uk", "DEF2")
	expectedErrorMessage := "Shortened URL already exists"

	_, err := fs.Create(longURLAlreadyShortened)
	if err.Error() != expectedErrorMessage {
		t.Errorf("Expected error message of '%s', instead received '%s'", expectedErrorMessage, err.Error())
	}

	// clean up
	clearTestData()
}

func TestItSuccessfullyCreatesAShortenedURLIfNoDataAlreadyExists(t *testing.T) {
	// clean up
	clearTestData()

	fs := getTestFsRepository()

	shortenedURL := shortenedurl.New("http://bbc.co.uk", "ABC1")

	result, err := fs.Create(shortenedURL)

	if err != nil {
		t.Errorf("Not expecting error, instead received '%s'", err.Error())
	}

	if result != shortenedURL {
		t.Errorf(
			"Expected identical ShortenedURL objects, instead received '%+v' and '%+v'",
			result,
			shortenedURL,
		)
	}
}

func TestItSuccessfullyCreatesAShortenedURLIfDifferentDataAlreadyExists(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	shortenedURL := shortenedurl.New("http://wikipedia.org.uk", "DEF2")

	result, err := fs.Create(shortenedURL)

	if err != nil {
		t.Errorf("Not expecting error, instead received '%s'", err.Error())
	}

	if result != shortenedURL {
		t.Errorf(
			"Expected identical ShortenedURL objects, instead received '%+v' and '%+v'",
			result,
			shortenedURL,
		)
	}

	// clean up
	clearTestData()
}

func TestItFailsToRetrieveByShortCodeThatDoesNotExist(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	_, err := fs.RetrieveByShortCode("DEF2")

	if err.Error() != "Shortened URL does not exist" {
		t.Errorf(
			"Expected error '%s', instead received '%s'",
			"Shortened URL does not exist",
			err.Error(),
		)
	}

	// clean up
	clearTestData()
}

func TestItSuccessfullyRetrievesByShortCode(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	shortenedURL, err := fs.RetrieveByShortCode("ABC1")

	if err != nil {
		t.Errorf("Not expecting error, instead received '%s'", err.Error())
	}

	if shortenedURL.GetLong() != "http://bbc.co.uk" {
		t.Errorf(
			"Expected long URL '%s', instead received '%s'",
			"http://bbc.co.uk",
			shortenedURL.GetLong(),
		)
	}

	if shortenedURL.GetShort() != "ABC1" {
		t.Errorf(
			"Expected shortcode '%s', instead received '%s'",
			"ABC1",
			shortenedURL.GetShort(),
		)
	}

	// clean up
	clearTestData()
}

func TestItFailsToRetrieveBySLongURLThatDoesNotExist(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	_, err := fs.RetrieveByLongURL("http://wikipedia.org")

	if err.Error() != "Shortened URL does not exist" {
		t.Errorf(
			"Expected error '%s', instead received '%s'",
			"Shortened URL does not exist",
			err.Error(),
		)
	}

	// clean up
	clearTestData()
}

func TestItSuccessfullyRetrievesByLongURL(t *testing.T) {
	// set expected data
	setTestData(`{"http://bbc.co.uk": "ABC1"}`)

	fs := getTestFsRepository()

	shortenedURL, err := fs.RetrieveByLongURL("http://bbc.co.uk")

	if err != nil {
		t.Errorf("Not expecting error, instead received '%s'", err.Error())
	}

	if shortenedURL.GetLong() != "http://bbc.co.uk" {
		t.Errorf(
			"Expected long URL '%s', instead received '%s'",
			"http://bbc.co.uk",
			shortenedURL.GetLong(),
		)
	}

	if shortenedURL.GetShort() != "ABC1" {
		t.Errorf(
			"Expected shortcode '%s', instead received '%s'",
			"ABC1",
			shortenedURL.GetShort(),
		)
	}

	// clean up
	clearTestData()
}

func TestItGetsPathToDBFile(t *testing.T) {
	fs := getTestFsRepository()
	path := getPathToDbFile(fs)

	expectedPath := fs.basePath + "/db.txt"

	if path != expectedPath {
		t.Errorf("Expected path '%s', instead received '%s'", expectedPath, path)
	}
}

func TestItSuccessfullyLoadsManifest(t *testing.T) {
	// set expected data
	setTestData(`{"hello": "world", "bonjour": "monde"}`)

	m := loadManifest(getTestDataPath())

	if len(m) != 2 {
		t.Errorf("Expected manifest length of %d, instead received %d", 2, len(m))
	}

	if m["hello"] != "world" {
		t.Errorf("Expected manifest value of '%s', instead received '%s'", "world", m["hello"])
	}

	if m["bonjour"] != "monde" {
		t.Errorf("Expected manifest value of '%s', instead received '%s'", "monde", m["bonjour"])
	}

	// clean up
	clearTestData()
}

func TestItSuccessfullySavesManifest(t *testing.T) {
	// set initial data
	setTestData(`{"hello": "world", "bonjour": "monde"}`)

	expectedMap := map[string]string{
		"goodbye":         "earth",
		"au revoir":       "terre",
		"auf wiedersehen": "erde",
	}

	result := saveManifest(getTestDataPath(), expectedMap)
	if result != true {
		t.Errorf("Expected save manifest to return true, instead returned '%+v'", result)
	}

	reloaded := loadManifest(getTestDataPath())

	if len(reloaded) != 3 {
		t.Errorf("Expected manifest length of %d, instead received %d", 2, len(reloaded))
	}

	if reloaded["goodbye"] != "earth" {
		t.Errorf("Expected manifest value of '%s', instead received '%s'", "earth", reloaded["goodbye"])
	}

	if reloaded["au revoir"] != "terre" {
		t.Errorf("Expected manifest value of '%s', instead received '%s'", "terre", reloaded["au revoir"])
	}

	if reloaded["auf wiedersehen"] != "erde" {
		t.Errorf("Expected manifest value of '%s', instead received '%s'", "erde", reloaded["auf wiedersehen"])
	}

	// clean up
	clearTestData()
}

func setTestData(data string) {
	clearTestData()
	ioutil.WriteFile(getTestDataPath(), []byte(data), 0644)
}

func clearTestData() {
	os.Remove(getTestDataPath())
}

func getTestDataPath() string {
	return getPathToDbFile(getTestFsRepository())
}

func getTestFsRepository() FileSystem {
	dir, _ := os.Getwd()

	return New(dir)
}
