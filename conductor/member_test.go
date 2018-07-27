package conductor

import (
	"context"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

func TestMemberAlreadyStopped(t *testing.T) {
	member := NewMember("", nil, provider.Request{}, nil)
	member.Stop()
	go member.Play()

	select {
	case err := <-member.DoneCh():
		if err != context.Canceled {
			t.Fatalf("got unexpected error instead of context canceled: %v", err)
		}
	}
}
