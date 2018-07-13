package conductor

import (
	"testing"
	"time"

	"github.com/tlmiller/disttrust/provider"
)

type testLease func() time.Time

func (t testLease) ID() string                            { return "testlease" }
func (t testLease) HasResponse() bool                     { return false }
func (t testLease) Response() (*provider.Response, error) { return nil, nil }
func (t testLease) Start() time.Time                      { return time.Now() }
func (t testLease) Till() time.Time                       { return t() }

// Tests calculating grace time with a zero duration. Expected result is that
// the grace should also be zero
func TestCalculateGraceWithZeroDuration(t *testing.T) {
	grace := calculateGrace(time.Duration(0))
	if grace != time.Duration(0) {
		t.Fatalf("expected zero grace duration for zero lease but got '%d'", grace)
	}
}

func TestCalculateGraceWithMinNanosecond(t *testing.T) {
	grace := calculateGrace(9 * time.Nanosecond)
	if grace != time.Duration(0) {
		t.Fatalf("expected zero grace duration for less than 10 nanoseconds but got '%d'", grace)
	}
}

func TestCalculateGraceRange(t *testing.T) {
	tests := []struct {
		lease time.Duration
		max   time.Duration
	}{
		{10 * time.Second, time.Second},
		{time.Second, 100 * time.Millisecond},
	}

	for _, test := range tests {
		grace := calculateGrace(test.lease)
		if grace < time.Duration(0) || grace > test.max {
			t.Fatalf("grace is outside of 80-90%% range with '%d' for a max of '%d'",
				grace, test.max)
		}
	}
}

func TestRenewWithinThreshold(t *testing.T) {
	r := NewRenewer(testLease(func() time.Time { return time.Now().Add(RenewThreshold) }))

	go r.Renew()

	renewChCall := false
	select {
	case <-r.RenewCh():
		renewChCall = true
	case err := <-r.DoneCh():
		if err != nil {
			t.Fatalf("go unexpected error from done channel: %v", err)
		}
		if renewChCall == false {
			t.Fatal("done channel called before renew channel")
		}
	}
}

func TestRenewWithStop(t *testing.T) {
	r := NewRenewer(testLease(func() time.Time { return time.Now().Add(time.Minute) }))

	go r.Renew()
	r.Stop()

	select {
	case <-r.RenewCh():
		t.Fatal("recieved renew channel call when not expected")
	case err := <-r.DoneCh():
		if err != nil {
			t.Fatalf("go unexpected error from done channel: %v", err)
		}
	}
}
