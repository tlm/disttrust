package config

import (
	"github.com/tlmiller/disttrust/action"
)

type Action struct {
	Command []string `json:"command"`
}

func ToAction(caction Action) (action.Action, error) {
	return action.CommandFromSlice(caction.Command)
}
