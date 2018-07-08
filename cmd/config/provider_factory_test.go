package config

import (
	"encoding/json"
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

func TestNoMappers(t *testing.T) {
	rawOpt := json.RawMessage(`{}`)
	provider, err := MakeProvider("test", rawOpt)
	if err == nil {
		t.Fatalf("MakeProvider should have errored with no mapper")
	}
	if provider != nil {
		t.Fatalf("non nil provider returned for failed MakeProvider")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	mapper := ProviderMapper(func(opt json.RawMessage) (provider.Provider, error) {
		return nil, nil
	})
	id := provider.Id("test")

	err := RegisterProvider(id, mapper)
	if err != nil {
		t.Errorf("recieved error for first register: %v", err)
	}

	err = RegisterProvider(id, mapper)
	if err == nil {
		t.Error("expected error for duplicate register")
	}
}
