package conductor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/dest"
	logu "github.com/tlmiller/disttrust/log"
	"github.com/tlmiller/disttrust/provider"
)

func ResponseHandle(dest dest.Dest, next LeaseHandler) LeaseHandler {
	h := func(ctx context.Context, lease provider.Lease) error {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		log := logu.GetLogger(ctx)
		log.Info("handling new lease response")

		if !lease.HasResponse() {
			log.Info("lease has no response, not sending to dest")
			return next.Handle(ctx, lease)
		}

		res, err := lease.Response()
		if err != nil {
			return errors.Wrap(err, "getting lease response for dest")
		}
		log.Info("sending new lease response to dest")
		err = dest.Send(res)
		if err != nil {
			return errors.Wrap(err, "sending response to dest")
		}
		return next.Handle(ctx, lease)
	}
	return LeaseHandlerFunc(h)
}
