package config

import (
	"testing"

	"github.com/tlmiller/disttrust/provider"
)

func TestProviderFactoryAddDoubleMapping(t *testing.T) {
	mapping := ProviderMapping{
		Config: func() interface{} { return nil },
		Mapper: func(_ interface{}) (provider.Provider, error) {
			return nil, nil
		},
	}
	pf := ProviderFactory{}

	if err := pf.AddMapping(provider.Id("test"), mapping); err != nil {
		t.Fatalf("unexpected error add provider mapper: %v", err)
	}

	if err := pf.AddMapping(provider.Id("test"), mapping); err == nil {
		t.Fatal("exepected error for adding double mapping on id test")
	}
}
