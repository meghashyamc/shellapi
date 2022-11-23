package models

type CmdRequest struct {
	ShellName string `json:"shell_name" validate:"max=128"`
	Password  string `json:"password" validate:"max=256"`
	Command   string `json:"command" validate:"required,min=1,max=256"`
}
