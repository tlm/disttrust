package conductor

import (
	"context"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type testAction struct {
	Fired     bool
	FiredFunc func(context.Context) error
}

func (a *testAction) Fire(ctx context.Context) error {
	a.Fired = true
	if a.FiredFunc != nil {
		return a.FiredFunc(ctx)
	}
	return nil
}

func TestActionHandleWithCanceledContext(t *testing.T) {
	nextHas := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		nextHas = true
		return nil
	})

	action := testAction{
		Fired: false,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ActionHandle(&action, next).Handle(ctx, &provider.DummyLease{})
	if err != nil {
		t.Fatalf("received unexpected error for no lease response: %v", err)
	}
	if nextHas {
		t.Fatalf("action handle called next when canceled context")
	}
	if action.Fired {
		t.Fatalf("exception action to not be fired for canceled context")
	}
}

func TestActionHandleWithNoResponse(t *testing.T) {
	nextHas := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		nextHas = true
		return nil
	})

	action := testAction{
		Fired: false,
	}

	err := ActionHandle(&action, next).Handle(context.Background(), &provider.DummyLease{})
	if err != nil {
		t.Fatalf("received unexpected error for no lease response: %v", err)
	}
	if !nextHas {
		t.Fatalf("action handle did not pass through to next")
	}
	if action.Fired {
		t.Fatalf("action fire when no lease response")
	}
}

func TestActionHandleFireForResponse(t *testing.T) {
	nextHas := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		nextHas = true
		return nil
	})

	action := testAction{
		Fired: false,
	}

	err := ActionHandle(&action, next).Handle(context.Background(),
		&provider.DummyLease{ResponseVal: &provider.Response{}})
	if err != nil {
		t.Fatalf("received unexpected error for no lease response: %v", err)
	}
	if !nextHas {
		t.Fatalf("action handle did not pass through to next")
	}
	if !action.Fired {
		t.Fatalf("action fire was not called")
	}
}

func TestActionHandleFireErrForDeadline(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	nextHas := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		nextHas = true
		return nil
	})

	action := testAction{
		Fired: false,
		FiredFunc: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	err := ActionHandle(&action, next).Handle(context.Background(),
		&provider.DummyLease{ResponseVal: &provider.Response{}})
	if err == nil {
		t.Fatalf("received no error for context deadline: %v", err)
	}
	if !action.Fired {
		t.Fatalf("action fire was not called")
	}
	if nextHas {
		t.Fatalf("action handle was called on failure")
	}
}
