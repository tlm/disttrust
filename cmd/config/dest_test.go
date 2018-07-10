package config

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUnknownDestId(t *testing.T) {
	tests := []string{
		"noexist",
		"test",
		"1234",
	}

	for _, test := range tests {
		_, err := MakeDest(test, json.RawMessage([]byte{}))
		if err == nil {
			t.Fatalf("expected unknown dest type error for %s", test)
		}

		expect := fmt.Sprintf("unknown dest type '%s'", test)
		if err.Error() != expect {
			t.Fatalf("unexpected error message for uknown dest type %s", test)
		}
	}
}
