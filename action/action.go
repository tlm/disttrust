package action

import (
	"context"
)

type Action interface {
	Fire(context.Context) error
}
