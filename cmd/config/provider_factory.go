package config

import (
	"encoding/json"
	"fmt"

	"github.com/tlmiller/disttrust/provider"
)

type ProviderMapper func(opt json.RawMessage) (provider.Provider, error)

var (
	providerFactories = make(map[provider.Id]ProviderMapper)
)

func MakeProvider(id string, opts json.RawMessage) (provider.Provider, error) {
	mapper, exists := providerFactories[provider.Id(id)]
	if exists == false {
		return nil, fmt.Errorf("provider mapper does not exists for id '%s'", id)
	}

	return mapper(opts)
}

func RegisterProvider(id provider.Id, mapper ProviderMapper) error {
	if _, exists := providerFactories[id]; exists {
		return fmt.Errorf("provider-mapper already registered for id '%s'", id)
	}
	providerFactories[id] = mapper
	return nil
}
