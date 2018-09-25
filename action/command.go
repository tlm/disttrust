package action

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
)

type Command struct {
	slice []string
}

func CommandFromSlice(slice []string) (*Command, error) {
	if len(slice) == 0 {
		return nil, errors.New("cannot make command from empty slice")
	}
	return &Command{
		slice: slice,
	}, nil
}

func (c *Command) Fire(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, c.slice[0], c.slice[1:]...)
	cmd.Stderr = nil
	cmd.Stdout = nil
	_, err := cmd.Output()

	if exErr, ok := err.(*exec.ExitError); err != nil && ok {
		return fmt.Errorf("command action: %v - %s", exErr, exErr.Stderr)
	} else if err != nil {
		return fmt.Errorf("command action: %v", exErr)
	}
	return nil
}
