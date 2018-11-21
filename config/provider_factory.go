package config

import (
	"fmt"

	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
	vaultconf "github.com/tlmiller/disttrust/vault/config"
)

type ProviderConfig func() interface{}

type ProviderMapper func(val interface{}) (provider.Provider, error)

type ProviderMapping struct {
	Config ProviderConfig
	Mapper ProviderMapper
}

type ProviderFactory map[provider.Id]ProviderMapping

var (
	DefaultProviderFactory = make(ProviderFactory)
)

func (p ProviderFactory) AddMapping(id provider.Id, m ProviderMapping) error {
	if _, exists := p[id]; exists {
		return fmt.Errorf("provider mapping already registered for id '%s'", id)
	}
	p[id] = m
	return nil
}

func init() {
	vaultMapping := ProviderMapping{
		Config: vaultconf.New,
		Mapper: vaultconf.Mapper,
	}
	if err := DefaultProviderFactory.AddMapping(provider.Id("vault"),
		vaultMapping); err != nil {
		panic(err)
	}
}

func (p ProviderMapping) MakeProvider(options map[string]interface{}) (
	provider.Provider, error) {
	conf := p.Config()
	if err := mapstructure.Decode(options, conf); err != nil {
		return nil, errors.Wrap(err, "making provider mapping config from options")
	}
	return p.Mapper(conf)
}
