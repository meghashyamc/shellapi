package shell

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const shellTimeout = 5 * time.Second

type output struct {
	Error      error
	StdOut     string
	StdErr     string
	Message    string
	StatusCode int
}

func ExecuteCmd(ctx context.Context, shellName, password string, command string, lg *log.Entry) *output {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	ctx, cancel := context.WithTimeout(ctx, shellTimeout)
	defer cancel()

	cmd := buildFullCommand(shellName, command)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = strings.NewReader(password)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	go func() {
		<-ctx.Done()
		p := cmd.Process
		if p == nil {
			return
		}
		killProcess(p, command, password, lg)
	}()
	err := cmd.Run()
	if err != nil {
		message, statusCode := getMessageAndStatusCodeForError(err)
		result := &output{Error: err, StdOut: stdout.String(), StdErr: stderr.String(), Message: message, StatusCode: statusCode}
		lg.WithFields(log.Fields{"err": err.Error(), "status": result.StatusCode}).Error("could not execute shell command")
		return result
	}
	return &output{Error: nil, StdOut: stdout.String(), StdErr: stderr.String(), Message: "command executed successfully", StatusCode: http.StatusOK}
}

func buildFullCommand(shellName string, command string) *exec.Cmd {

	commandsList := getCleanupCommandsList(command)

	if len(shellName) > 0 {
		commandsListWithShell := []string{}
		commandsListWithShell = append(commandsListWithShell, shellName, "-c")
		commandsListWithShell = append(commandsListWithShell, strings.Join(commandsList, " "))
		return exec.Command(commandsListWithShell[0], commandsListWithShell[1:]...)
	}

	if len(commandsList) == 1 {
		return exec.Command(commandsList[0])
	}
	return exec.Command(commandsList[0], commandsList[1:]...)

}

func getCleanupCommandsList(command string) []string {

	commandsList := strings.Split(command, " ")

	cleanedUpCommands := []string{}

	for _, singleCommand := range commandsList {

		cleanedUpCommands = append(cleanedUpCommands, strings.TrimSpace(singleCommand))
	}

	return cleanedUpCommands
}

func getMessageAndStatusCodeForError(err error) (string, int) {

	if isNotFoundErr(err) {

		return errCommandNotFound, http.StatusNotFound
	}
	if isTimeOutErr(err) {
		return errCommandTimedOut, http.StatusRequestTimeout
	}
	return errInternal, http.StatusInternalServerError

}

func killProcess(p *os.Process, command, password string, lg *log.Entry) {
	errProcessKillFailed := errors.New("could not kill process after trying to execute shell command till timeout")

	if !strings.Contains(command, "sudo") {
		if err := syscall.Kill(-p.Pid, syscall.SIGKILL); err != nil {
			lg.WithFields(log.Fields{"err": err.Error()}).Error(errProcessKillFailed.Error())
			return
		}
	}

	if _, err := exec.Command("bash", "-c", fmt.Sprintf("echo %s | sudo -S kill -%d -%d", password, syscall.SIGKILL, p.Pid)).Output(); err != nil {

		lg.WithFields(log.Fields{"err": err.Error()}).Error(errProcessKillFailed.Error())

	}

}

func isNotFoundErr(err error) bool {

	return err != nil && strings.Contains(err.Error(), "file not found")
}

func isTimeOutErr(err error) bool {
	return err != nil && strings.Contains(err.Error(), "killed")

}
