package conductor

import (
	"testing"
)

// Tests the correct bahaviour is shown when stop is called on a member that has
// never had a subsequent Play() call. This call should result in an error
// saying the member is already stopped
func TestMemberNoInitStop(t *testing.T) {
	member := Member{}
	err := member.Stop()
	if err == nil {
		t.Fatal("expected err value for Stop() call on non played member")
	}

	if err.Error() != "member already stopped" {
		t.Fatal("stop call did not produce expected error message")
	}
}
