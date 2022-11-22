package models

type CmdRequest struct {
	Command   string   `json:"command" validate:"required,alpha,min=1,max=128"`
	Arguments []string `json:"arguments" validate:"dive,required,min=1,max=128"`
}
