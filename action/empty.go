package action

import (
	"context"
)

type Empty struct {
}

func (e *Empty) Fire(_ context.Context) error {
	return nil
}
