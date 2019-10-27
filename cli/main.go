package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type handler struct {
	command string
	param   string
}

const apiBaseURL = "http://localhost:8080"

func main() {
	cliCmd := "go run cli/main.go"

	handler := newHandler(os.Args)

	commands := map[string]func(s string) error{
		"shorten":  commandShorten,
		"redirect": commandRedirect,
	}

	if commands[handler.command] != nil {
		err := commands[handler.command](handler.param)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}

		return
	}

	// fallback (no command supplied)
	fmt.Println("Usage:")
	fmt.Printf("%s shorten <url>           Shorten a long URL\n", cliCmd)
	fmt.Printf("%s redirect <shortcode>    Redirect a shortcode to original URL\n", cliCmd)
}

func newHandler(args []string) handler {
	if len(args) == 1 {
		return handler{}
	}

	command := args[1]
	var param string

	if len(args) > 2 {
		param = args[2]
	}

	return handler{
		command: command,
		param:   param,
	}
}

func commandShorten(param string) error {
	if param == "" {
		return errors.New("Please supply a URL to shorten")
	}

	// build request payload and make request
	requestPayload := fmt.Sprintf("{\"url\": \"%s\"}", param)
	resp, err := http.Post(
		fmt.Sprintf("%s/api/shorten", apiBaseURL),
		"application/json",
		strings.NewReader(requestPayload),
	)
	if err != nil {
		return errors.New(err.Error())
	}

	// read response body
	body := []byte{}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New(err.Error())
	}

	// parse body as json
	parsed := map[string]interface{}{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return fmt.Errorf("Status %d, %s", resp.StatusCode, err.Error())
	}

	// check if response body indicates a success
	if parsed["status"].(string) != "ok" {
		message := parsed["data"].(map[string]interface{})["message"]
		return errors.New(message.(string))
	}

	// print short URL
	shortURL := parsed["data"].(map[string]interface{})["shortURL"]
	fmt.Println(shortURL.(string))

	return nil
}

func commandRedirect(param string) error {
	// @TODO Implement
	return nil
}
