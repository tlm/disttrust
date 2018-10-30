package dest

import (
	"errors"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

type mockDest struct {
	counter *int
	fail    bool
}

func (m *mockDest) Send(_ *provider.Response) error {
	(*m.counter)++
	if m.fail {
		return errors.New("I was told to fail")
	}
	return nil
}

func TestAggregateCallsAllDests(t *testing.T) {
	var counter int
	agg := NewAggregate(&mockDest{&counter, false}, &mockDest{&counter, false},
		&mockDest{&counter, false})
	err := agg.Send(nil)
	if err != nil {
		t.Fatalf("unexpected error for aggregate send: %v", err)
	}
	if counter != 3 {
		t.Fatalf("aggregate did not call all dests, expected 3 got %d", counter)
	}
}

func TestAggregateCallDestFailure(t *testing.T) {
	var counter int
	agg := NewAggregate(&mockDest{&counter, false}, &mockDest{&counter, false},
		&mockDest{&counter, true}, &mockDest{&counter, false},
		&mockDest{&counter, false})
	err := agg.Send(nil)
	if err == nil {
		t.Fatal("expected non nill err for aggregate failure send")
	}
	if counter != 3 {
		t.Fatalf("aggregate did not call all dests, expected 3 got %d", counter)
	}
}
