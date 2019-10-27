# HTTP URL Shortener 🔗

Small HTTP service that accepts a long URL to be shortened, and generates a shortcode that redirects to its origin.

## Requirements

* Golang 1.13

## Getting Started

From the project root directory:

```
go build -o api-bin -i ./api
./api-bin
```

This will launch a HTTP server for the URL Shortener service, listening locally on port `8080`.

## Usage

### API

To shorten a URL, make the following request:

```
curl -X POST \
  http://localhost:8080/api/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"http://bbc.co.uk"}'
```

This will return a response payload - e.g.:

```
{
  "status": "ok",
    "data": {
      "shortURL": "http://localhost:8080/ABC1"
  }
}
```

Make the following request (or visit this URL in your browser to be redirected
to the original source URL):

```
curl -X GET \
  http://localhost:8080/ABC1
```

This will return a redirect response - e.g.:

```
HTTP/1.1 301 Moved Permanently
Content-Type: application/json
Location: http://bbc.co.uk
```

### Command Line Interface

Whilst the API is running, you can issue the following commands
via a second terminal from the project root directory:

```
go run cli/main.go shorten <url>

go run cli/main.go redirect <shortcode>
```

Example:

```
go run cli/main.go   // cli usage instructions

go run cli/main.go shorten http://bbc.co.uk   // output: http://localhost:8080/ABC1

go run cli/main.go redirect ABC1              // launches http://bbc.co.uk in the default web browser
```

## Tests

To run the full project testsuite

```
go test ./...
```

## Roadmap

The API is backed by a File System data store, facilitated by the `FileSystem` type within the `shortenedurlfilesystemrepository` package.

This package implements the `RepositoryInterface` type within the `repositoryinterface` package, which the API's handler methods accept as a dependency.

These interface methods can be implemented on additional repository types that facilitate alternative data stores or caches (MySQL, Redis etc.)

Env/config could then be used to determine which repository to instantiate within the API's `main()` method, before passing to a handler method. 
