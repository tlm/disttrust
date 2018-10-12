package provider

import (
	"testing"
)

func TestDupeNameProviders(t *testing.T) {
	store := NewStore()
	err := store.Store("test", nil)
	if err != nil {
		t.Fatalf("receeved error for provider store: %v", err)
	}

	err = store.Store("test", nil)
	if err == nil {
		t.Fatal("should have received error for duplicate provider name")
	}
}

func TestProvidersRemoval(t *testing.T) {
	store := NewStore()
	p := &DummyProvider{}
	err := store.Store("test1", p)
	if err != nil {
		t.Fatalf("receeved error for provider store: %v", err)
	}

	pf, err := store.Fetch("test1")
	if err != nil {
		t.Fatalf("receeved error for provider fetch: %v", err)
	}

	if p != pf {
		t.Fatal("provider store returned does not match that stored")
	}

	store.Remove("test1")
	pf, err = store.Fetch("test1")
	if err == nil {
		t.Fatal("provider store fetch should have failed for removed povider")
	}
}
