package config

import (
	"testing"

	"github.com/tlmiller/disttrust/dest"
)

func TestGetDestWithEmptyValues(t *testing.T) {
	mapDest, err := GetDest(&DestConfig{})
	if err != nil {
		t.Fatalf("unexpected error while getting dest with empty values: %v", err)
	}

	_, ok := mapDest.(*dest.Empty)
	if !ok {
		t.Fatalf("unexpected dest type for empty values")
	}
}
