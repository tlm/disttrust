package conductor

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"

	"github.com/tlmiller/disttrust/action"
	logu "github.com/tlmiller/disttrust/log"
	"github.com/tlmiller/disttrust/provider"
)

var (
	ActionTimeout = 30 * time.Second
)

func ActionHandle(action action.Action, next LeaseHandler) LeaseHandler {
	h := func(ctx context.Context, lease provider.Lease) error {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		log := logu.GetLogger(ctx)
		log.Info("handling new lease action")

		if !lease.HasResponse() {
			log.Info("lease has no response, not firing action")
			return next.Handle(ctx, lease)
		}

		log.WithFields(logrus.Fields{
			"timeout": ActionTimeout / time.Second,
		})
		log.Info("firing new lease action")
		fireCtx, cancel := context.WithTimeout(ctx, ActionTimeout)
		defer cancel()

		err := action.Fire(fireCtx)
		if err != nil {
			return errors.Wrap(err, "firing action for new lease response")
		}
		return next.Handle(ctx, lease)
	}
	return LeaseHandlerFunc(h)
}
