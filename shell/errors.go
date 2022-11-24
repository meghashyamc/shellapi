package shell

import (
	"errors"
	"net/http"
	"strings"
)

var (
	ErrNoCommand          = errors.New("no command was sent")
	ErrInvalidSudoCommand = errors.New("command with sudo not sent in the expected format")
)

func getMessageAndStatusCodeForError(err error) (string, int) {

	if isNotFoundErr(err) {

		return errCommandNotFound, http.StatusNotFound
	}
	if isTimeOutErr(err) {
		return errCommandTimedOut, http.StatusRequestTimeout
	}
	return errExitCodeNotZero, http.StatusBadRequest

}

func isNoSuchProcessErr(err error) bool {

	return err != nil && strings.Contains(err.Error(), "no such process")
}
func isNotFoundErr(err error) bool {

	return err != nil && strings.Contains(err.Error(), "file not found")
}

func isTimeOutErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "killed")

}
