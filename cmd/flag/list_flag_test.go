package flag

import (
	"fmt"
	"testing"
)

func TestListFlagSetMany(t *testing.T) {
	list := ListFlag{}
	for i := 0; i < 256; i++ {
		set := fmt.Sprintf("another set %d", i)
		list.Set(set)
		if list[i] != set {
			t.Fatal("list set on index did not result in expected value")
		}
	}
}
