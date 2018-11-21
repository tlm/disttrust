package config

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/provider"
)

func TestGetProviderCalls(t *testing.T) {
	called := false
	pf := ProviderFactory{}

	mapping := ProviderMapping{
		Config: func() interface{} { return &struct{}{} },
		Mapper: func(_ interface{}) (provider.Provider, error) {
			called = true
			return nil, nil
		},
	}
	err := pf.AddMapping(provider.Id("test"), mapping)
	if err != nil {
		t.Fatalf("unexpected error adding new provider factory mapping: %v", err)
	}

	providers := []map[string]interface{}{
		map[string]interface{}{
			"id":   "test",
			"name": "test",
		},
	}
	v := viper.New()
	v.Set(Providers, providers)
	if _, err := GetProviderStoreWithFactory(v, provider.NewStore(), pf); err != nil {
		t.Fatalf("unexpected error while getting mapped provider: %v", err)
	}

	if !called {
		t.Fatal("provider test mapper was not called")
	}
}

func TestGetProviderReuse(t *testing.T) {
	called := false
	pf := ProviderFactory{}
	mapping := ProviderMapping{
		Config: func() interface{} { return &struct{}{} },
		Mapper: func(_ interface{}) (provider.Provider, error) {
			called = true
			return nil, nil
		},
	}

	err := pf.AddMapping(provider.Id("test"), mapping)
	if err != nil {
		t.Fatalf("unexpected error adding new provider factory mapping: %v", err)
	}

	p := &ProviderConfigWrap{
		Options: map[string]interface{}{
			"test": "val",
		},
	}
	store := provider.NewStore()
	store.Store("test", p)

	providers := []map[string]interface{}{
		map[string]interface{}{
			"id":   "test",
			"name": "test",
			"options": map[string]interface{}{
				"test": "val",
			},
		},
	}
	v := viper.New()
	v.Set(Providers, providers)

	nstore, err := GetProviderStoreWithFactory(v, store, pf)
	if err != nil {
		t.Fatalf("unexpected error while getting mapped provider: %v", err)
	}

	if _, exists := nstore.Fetch("test"); !exists {
		t.Fatal("provider test does not exist in returned store")
	}

	if called {
		t.Fatal("provider mapping function was called when it should have been reused")
	}
}

func TestGetProviderDoesNotReuse(t *testing.T) {
	called := false
	pf := ProviderFactory{}
	mapping := ProviderMapping{
		Config: func() interface{} { return &struct{}{} },
		Mapper: func(_ interface{}) (provider.Provider, error) {
			called = true
			return nil, nil
		},
	}

	err := pf.AddMapping(provider.Id("test"), mapping)
	if err != nil {
		t.Fatalf("unexpected error adding new provider factory mapping: %v", err)
	}

	p := &ProviderConfigWrap{
		Options: map[string]interface{}{
			"test": "val",
		},
	}
	store := provider.NewStore()
	store.Store("test", p)

	providers := []map[string]interface{}{
		map[string]interface{}{
			"id":   "test",
			"name": "test",
			"options": map[string]interface{}{
				"test": "val1",
			},
		},
	}
	v := viper.New()
	v.Set(Providers, providers)

	nstore, err := GetProviderStoreWithFactory(v, store, pf)
	if err != nil {
		t.Fatalf("unexpected error while getting mapped provider: %v", err)
	}

	if _, exists := nstore.Fetch("test"); !exists {
		t.Fatal("provider test does not exist in returned store")
	}

	if !called {
		t.Fatal("provider mapping function was called when it should have been reused")
	}
}
