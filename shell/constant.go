package shell

import "net/http"

const (
	errCommandNotFound = "could not identify command"
	errCommandTimedOut = "timed out when trying to execute command"
)

var errInternal = http.StatusText(http.StatusInternalServerError)
