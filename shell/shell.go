package shell

import (
	"bytes"
	"context"
	"net/http"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const shellTimeout = 5 * time.Second

type output struct {
	Error      error
	StdOut     string
	StdErr     string
	StatusCode int
}

func ExecuteCmd(cmdToExecute string, arguments []string) *output {
	ctx, cancel := context.WithTimeout(context.Background(), shellTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdToExecute, arguments...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		result := &output{Error: err, StdOut: stdout.String(), StdErr: stderr.String(), StatusCode: getStatusCodeForError(err)}

		log.WithFields(log.Fields{"err": err.Error(), "status": result.StatusCode}).Error("could not execute shell command")

		return result
	}
	return &output{Error: nil, StdOut: stdout.String(), StdErr: stderr.String(), StatusCode: http.StatusOK}
}

func getStatusCodeForError(err error) int {

	if isNotFoundErr(err) {
		return http.StatusNotFound
	}
	if isTimeOutErr(err) {
		return http.StatusRequestTimeout
	}
	return http.StatusInternalServerError

}
func isNotFoundErr(err error) bool {

	return err != nil && strings.Contains(err.Error(), "file not found")
}

func isTimeOutErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "killed")

}
