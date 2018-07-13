package conductor

import (
	"context"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type testDest struct {
	called bool
}

func (d *testDest) Send(_ *provider.Response) error {
	d.called = true
	return nil
}

func TestResponseHandleCtxDone(t *testing.T) {
	passed := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		passed = true
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ResponseHandle(&testDest{}, next).Handle(ctx, &provider.DummyLease{})
	if err != nil {
		t.Fatalf("recieved unexpected error: %v", err)
	}
	if passed {
		t.Fatalf("lease passed through to next when context done")
	}
}

func TestResponseHandleNext(t *testing.T) {
	passed := false
	next := LeaseHandlerFunc(func(_ context.Context, _ provider.Lease) error {
		passed = true
		return nil
	})

	dest := testDest{}
	err := ResponseHandle(&dest, next).Handle(context.Background(),
		&provider.DummyLease{
			ResponseVal: &provider.Response{},
		})

	if err != nil {
		t.Fatalf("recieved unexpected error: %v", err)
	}
	if !dest.called {
		t.Fatalf("dest was not called by response handler")
	}
	if !passed {
		t.Fatalf("response handle did not call next")
	}
}
