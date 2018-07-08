package config

import (
	"encoding/json"
	"testing"
)

func TestBadValues(t *testing.T) {
	rawOpt := json.RawMessage(`{"address": 1234}`)

	_, err := vaultMapper(rawOpt)
	if err == nil {
		t.Error("expected error for no vault options")
	}
}
