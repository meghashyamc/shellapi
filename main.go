package main

import (
	"github.com/meghashyamc/shellapi/api"
)

func main() {

	api.LogSetup()
	api.NewHTTPListener().Listen()

}
