package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meghashyamc/shellapi/models"
	"github.com/meghashyamc/shellapi/shell"

	log "github.com/sirupsen/logrus"
)

func (l *HTTPListener) homeHandler(w http.ResponseWriter, r *http.Request) {
	lg := logger(mux.CurrentRoute(r).GetName(), r.Context())
	lg.Info()

	writeResponse(w, http.StatusOK, true, "Successfully reached API home", nil, nil)

	return
}

func (l *HTTPListener) cmdHandler(w http.ResponseWriter, r *http.Request) {

	lg := logger(mux.CurrentRoute(r).GetName(), r.Context())
	lg.Info()
	cmdRequest := models.CmdRequest{}
	if err := json.NewDecoder(r.Body).Decode(&cmdRequest); err != nil {
		lg.WithFields(log.Fields{"err": err.Error()}).Error(errCouldNotUnmarshal)
		writeResponse(w, http.StatusBadRequest, false, "", []string{errCouldNotUnmarshal}, nil)
		return
	}

	defer r.Body.Close()

	if err := l.validate.Struct(cmdRequest); err != nil {
		lg.WithFields(log.Fields{"err": err.Error()}).Info(errRequestValidationFailed)
		writeResponse(w, http.StatusBadRequest, false, errRequestValidationFailed, getValidationErrors(err), nil)
		return
	}

	commandsList, err := shell.ValidateCmd(cmdRequest.Command, lg)
	if err != nil {
		lg.WithFields(log.Fields{"err": err.Error()}).Info(errRequestValidationFailed)
		writeResponse(w, http.StatusBadRequest, false, errRequestValidationFailed, []string{err.Error()}, nil)
		return
	}
	result := shell.ExecuteCmd(r.Context(), cmdRequest.ShellName, cmdRequest.Password, commandsList, lg)
	if result.Error != nil {
		writeResponse(w, result.StatusCode, false, result.Message, []string{result.Error.Error()}, cmdResponse{StdOut: result.StdOut, StdErr: result.StdErr})
		return
	}
	writeResponse(w, http.StatusOK, true, result.Message, []string{}, cmdResponse{StdOut: result.StdOut, StdErr: result.StdErr})
	return
}
