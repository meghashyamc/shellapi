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

func ValidateCmd(command string, lg *log.Entry) ([]string, error) {

	commandsList := getCleanedUpCommandsList(command)

	if len(commandsList) == 0 {
		return nil, ErrNoCommand
	}

	if err := isSudoCommandValid(commandsList); err != nil {
		return nil, err
	}
	return commandsList, nil

}

func ExecuteCmd(ctx context.Context, shellName, password string, commandsList []string, lg *log.Entry) *output {
	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	ctx, cancel := context.WithTimeout(ctx, shellTimeout)
	defer cancel()

	cmd := buildFullCommand(shellName, commandsList)

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
		killProcess(p, commandsList[0], password, lg)
	}()
	err := cmd.Run()
	if err != nil {
		message, statusCode := getMessageAndStatusCodeForError(err)
		result := &output{Error: err, StdOut: stdout.String(), StdErr: stderr.String(), Message: message, StatusCode: statusCode}
		lg.WithFields(log.Fields{"err": err.Error(), "status": result.StatusCode}).Error(message)
		return result
	}
	return &output{Error: nil, StdOut: stdout.String(), StdErr: stderr.String(), Message: "command executed successfully", StatusCode: http.StatusOK}
}

func buildFullCommand(shellName string, commandsList []string) *exec.Cmd {

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

func isSudoCommandValid(commandsList []string) error {

	if commandsList[0] != sudoCommand {
		return nil
	}
	if len(commandsList) < 3 {
		return ErrInvalidSudoCommand
	}
	if commandsList[1] != "-S" {
		return ErrInvalidSudoCommand
	}
	return nil

}

func getCleanedUpCommandsList(command string) []string {

	commandsList := strings.Split(command, " ")

	cleanedUpCommands := []string{}

	for _, singleCommand := range commandsList {

		cleanedUpCommands = append(cleanedUpCommands, strings.TrimSpace(singleCommand))
	}

	return cleanedUpCommands
}

func killProcess(p *os.Process, firstCommand, password string, lg *log.Entry) {
	errProcessKillFailed := errors.New("could not kill process after trying to execute shell command till timeout")

	if _, err := os.FindProcess(int(-p.Pid)); err != nil {
		return
	}

	if firstCommand != sudoCommand {
		if err := syscall.Kill(-p.Pid, syscall.SIGKILL); err != nil && !isNoSuchProcessErr(err) {
			lg.WithFields(log.Fields{"err": err.Error()}).Error(errProcessKillFailed.Error())
		}
		return
	}
	if _, err := exec.Command("bash", "-c", fmt.Sprintf("echo %s | sudo -S kill -%d -%d", password, syscall.SIGKILL, p.Pid)).Output(); err != nil && !isNoSuchProcessErr(err) {

		lg.WithFields(log.Fields{"err": err.Error()}).Error(errProcessKillFailed.Error())

	}

}
