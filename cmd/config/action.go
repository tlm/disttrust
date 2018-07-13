package config

import (
	"github.com/tlmiller/disttrust/action"
)

func MakeAction(caction Action) (action.Action, error) {
	return action.CommandFromSlice(caction.Command)
}
