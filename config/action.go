package config

import (
	"github.com/tlmiller/disttrust/action"
)

type ActionConfig struct {
	Command []string
}

func GetAction(a *ActionConfig) (action.Action, error) {
	if len(a.Command) != 0 {
		return action.CommandFromSlice(a.Command)
	}
	return &action.Empty{}, nil
}
