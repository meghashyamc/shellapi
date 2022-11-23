package main

import (
	"os"

	"github.com/meghashyamc/shellapi/api"
)

func main() {

	api.LogSetup()
	listener, err := api.NewHTTPListener()
	if err != nil {
		os.Exit(1)
	}

	listener.Listen()

}
