package api

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type cmdResponse struct {
	StdOut string `json:"stdout"`
	StdErr string `json:"stderr"`
}

func writeResponse(w http.ResponseWriter, statusCode int, success bool, message string, errors []string, data interface{}) {

	jsonBytes, err := json.Marshal(response{Success: success, Message: message, Errors: errors, Data: data})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not marshal response")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(jsonBytes)
	return

}
