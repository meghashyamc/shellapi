package api

import (
	"encoding/json"
	"net/http"

	"github.com/meghashyamc/shellapi/models"
	"github.com/meghashyamc/shellapi/shell"

	log "github.com/sirupsen/logrus"
)

func (l *HTTPListener) homeHandler(w http.ResponseWriter, r *http.Request) {

	writeResponse(w, http.StatusOK, true, "Successfully reached API home", nil, nil)

	return
}

func (l *HTTPListener) cmdHandler(w http.ResponseWriter, r *http.Request) {

	log.Info("api/cmd request received")
	cmdRequest := models.CmdRequest{}
	if err := json.NewDecoder(r.Body).Decode(&cmdRequest); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error(errCouldNotUnmarshal)
		writeResponse(w, http.StatusBadRequest, false, "", []string{errCouldNotUnmarshal}, nil)
		return
	}

	defer r.Body.Close()

	if err := l.validate.Struct(cmdRequest); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Info(errRequestValidationFailed)
		writeResponse(w, http.StatusBadRequest, false, errRequestValidationFailed, getValidationErrors(err), nil)
		return
	}

	result := shell.ExecuteCmd(cmdRequest.Command, cmdRequest.Arguments)
	if result.Error != nil {
		writeResponse(w, result.StatusCode, false, "command execution failed", []string{result.Error.Error()}, cmdResponse{StdOut: result.StdOut, StdErr: result.StdErr})
		return
	}
	writeResponse(w, http.StatusOK, true, "command executed successfully", []string{}, cmdResponse{StdOut: result.StdOut, StdErr: result.StdErr})
	return
}
