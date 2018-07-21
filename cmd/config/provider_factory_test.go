package config

import (
	"encoding/json"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

func TestNoMappers(t *testing.T) {
	rawOpt := json.RawMessage(`{}`)
	provider, err := ToProvider("test", rawOpt)
	if err == nil {
		t.Fatalf("ToProvider should have errored with no mapper")
	}
	if provider != nil {
		t.Fatalf("non nil provider returned for failed ToProvider")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	mapper := ProviderMapper(func(opt json.RawMessage) (provider.Provider, error) {
		return nil, nil
	})
	id := provider.Id("test")

	err := MapProvider(id, mapper)
	if err != nil {
		t.Errorf("received error for first mapping: %v", err)
	}

	err = MapProvider(id, mapper)
	if err == nil {
		t.Error("expected error for duplicate map")
	}
}
