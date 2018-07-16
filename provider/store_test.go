package provider

import (
	"testing"
)

func TestDupeNameProviders(t *testing.T) {
	err := Store("test", nil)
	if err != nil {
		t.Fatalf("receeved error for provider store: %v", err)
	}

	err = Store("test", nil)
	if err == nil {
		t.Fatal("should have received error for duplicate provider name")
	}
}
