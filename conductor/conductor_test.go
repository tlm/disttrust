package conductor

import (
	"errors"
	"testing"
)

type testMember struct {
	doneCh    chan error
	playCount int
	playFunc  func() error
}

func (t *testMember) DoneCh() <-chan error {
	return t.doneCh
}

func (t *testMember) Play() {
	t.playCount++
	if t.playFunc == nil {
		t.doneCh <- nil
		return
	}
	t.doneCh <- t.playFunc()
}

func (t *testMember) Stop() {
}

func TestConductorMemberDone(t *testing.T) {
	member := testMember{
		doneCh: make(chan error),
	}

	c := NewConductor().AddMember(&member)
	err := c.Watch()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if member.playCount != 1 {
		t.Fatalf("expected")
	}
}

func TestConductorMaxRetries(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	member := testMember{
		doneCh:   make(chan error),
		playFunc: func() error { return errors.New("test error") },
	}

	c := NewConductor().AddMember(&member)
	err := c.Watch()
	if err == nil {
		t.Fatalf("expected error for member consuming the maximum number of retries")
	}
	if member.playCount != MaxRetry+1 {
		t.Fatalf("expected member to be played %d times but got %d",
			MaxRetry+1, member.playCount)
	}
}
