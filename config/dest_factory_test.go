package config

import (
	"testing"

	"github.com/tlmiller/disttrust/dest"
)

func TestDestFactoryDoubleMapping(t *testing.T) {
	mapping := DestMapping{
		Config: func() interface{} { return nil },
		Mapper: func(_ interface{}) (dest.Dest, error) {
			return nil, nil
		},
	}
	factory := DestFactory{}
	if err := factory.AddMapping("test", mapping); err != nil {
		t.Fatal("unexpected error when adding test dest mapping to factory")
	}

	if err := factory.AddMapping("test", mapping); err == nil {
		t.Fatal("expected error when adding double mapping to dest factory")
	}
}
