package conductor

import (
	"testing"
)

// Tests the correct bahaviour is shown when stop is called on a member that has
// never had a subsequent Play() call. This call should result in an error
// saying the member is already stopped
func TestMemberNoInitStop(t *testing.T) {
	//member := Member{}
	//	member.Stop()
}
