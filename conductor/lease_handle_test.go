package conductor

import (
	"context"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

// Testing that NewLeaseHandler calls the next func in the chain
func TestNewLeaseHandlePass(t *testing.T) {
	passed := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		passed = true
		return nil
	})

	err := NewLeaseHandle(next).Handle(context.Background(), &provider.DummyLease{})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
	if !passed {
		t.Fatal("new lease handle did not pass through to next")
	}
}

func TestNewLeaseContextDone(t *testing.T) {
	passed := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		passed = true
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := NewLeaseHandle(next).Handle(ctx, &provider.DummyLease{})
	if err != nil {
		t.Fatalf("received unexpected error: %v", err)
	}
	if passed {
		t.Fatal("new lease handle passed through to next when context done")
	}
}
