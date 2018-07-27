package conductor

import (
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

func TestDummy(t *testing.T) {
}
