package conductor

import (
	"context"

	"github.com/tlmiller/disttrust/action"
	"github.com/tlmiller/disttrust/dest"
	"github.com/tlmiller/disttrust/provider"
)

type LeaseHandler interface {
	Handle(context.Context, provider.Lease) error
}

type LeaseHandlerFunc func(context.Context, provider.Lease) error

var (
	NoOpLeaseHandle = LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		return nil
	})
)

func DefaultLeaseHandle(rdest dest.Dest, ract action.Action) LeaseHandler {
	return NewLeaseHandle(ResponseHandle(rdest, ActionHandle(ract, NoOpLeaseHandle)))
}

func (l LeaseHandlerFunc) Handle(ctx context.Context, lease provider.Lease) error {
	return l(ctx, lease)
}
