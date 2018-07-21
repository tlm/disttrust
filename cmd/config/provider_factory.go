package config

import (
	"encoding/json"
	"fmt"

	"github.com/tlmiller/disttrust/provider"
)

type ProviderMapper func(opt json.RawMessage) (provider.Provider, error)

var (
	providerMappings = make(map[provider.Id]ProviderMapper)
)

func MapProvider(id provider.Id, mapper ProviderMapper) error {
	if _, exists := providerMappings[id]; exists {
		return fmt.Errorf("provider mapping already registered for id '%s'", id)
	}
	providerMappings[id] = mapper
	return nil
}

func ToProvider(id string, opts json.RawMessage) (provider.Provider, error) {
	mapper, exists := providerMappings[provider.Id(id)]
	if exists == false {
		return nil, fmt.Errorf("provider mapper does not exists for id '%s'", id)
	}

	return mapper(opts)
}
