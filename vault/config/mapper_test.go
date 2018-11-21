package config

import (
	"testing"

	"github.com/spf13/viper"
)

func TestBadMapperValues(t *testing.T) {
	v := viper.New()
	r := New()
	if err := v.Unmarshal(r); err != nil {
		t.Fatalf("unexpected error unmarshalling vault config: %v", err)
	}
	_, err := Mapper(r)
	if err == nil {
		t.Error("expected error for no vault options")
	}

	v.Set("authMethod", "noexist")
	r = New()
	if err := v.Unmarshal(r); err != nil {
		t.Fatalf("unexpected error unmarshalling vault config: %v", err)
	}
	_, err = Mapper(r)
	if err == nil {
		t.Error("expected error for bas auth method")
	}
}
