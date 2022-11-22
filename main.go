package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/meghashyamc/shellapi/api"
	log "github.com/sirupsen/logrus"
)

func main() {

	godotenv.Load()
	log.SetFormatter(&log.JSONFormatter{})
	listener, err := api.NewHTTPListener()
	if err != nil {
		os.Exit(1)
	}
	listener.Listen()

}
