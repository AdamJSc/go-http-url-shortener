package main

import (
	"fmt"
	"os"
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
	// @TODO Implement
	return nil
}

func commandRedirect(param string) error {
	// @TODO Implement
	return nil
}
