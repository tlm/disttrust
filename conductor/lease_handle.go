package conductor

import (
	"context"

	logu "github.com/tlmiller/disttrust/log"
	"github.com/tlmiller/disttrust/provider"
)

func NewLeaseHandle(next LeaseHandler) LeaseHandler {
	h := func(ctx context.Context, lease provider.Lease) error {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		log := logu.GetLogger(ctx)
		log.Info("handling new lease")
		return next.Handle(ctx, lease)
	}
	return LeaseHandlerFunc(h)
}
