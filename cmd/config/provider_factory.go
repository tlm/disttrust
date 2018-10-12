package config

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"

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

	p, err := mapper(opts)
	return &ProviderConfigWrap{
		config: opts,
		p:      p,
	}, err
}

func ToProviderOnUpdate(id string, opts json.RawMessage, ex provider.Provider) (provider.Provider, error) {
	if cpwrap, ok := ex.(*ProviderConfigWrap); ok {
		var j1, j2 interface{}
		if err := json.Unmarshal(opts, &j1); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cpwrap.config, &j2); err != nil {
			return nil, err
		}
		if reflect.DeepEqual(j1, j2) {
			return ex, nil
		}
	}

	mapper, exists := providerMappings[provider.Id(id)]
	if exists == false {
		return nil, fmt.Errorf("provider mapper does not exists for id '%s'", id)
	}

	p, err := mapper(opts)
	return &ProviderConfigWrap{
		config: opts,
		p:      p,
	}, err
}

func ToProviderStore(cnfProviders []Provider, store *provider.Store) (*provider.Store, error) {
	nstore := provider.NewStore()
	for _, cnfProvider := range cnfProviders {
		if len(cnfProvider.Name) == 0 {
			return nstore, errors.New("undefined provider name")
		}

		var err error
		var genProvider provider.Provider
		if p, exists := store.Fetch(cnfProvider.Name); exists {
			genProvider, err = ToProviderOnUpdate(cnfProvider.Name,
				cnfProvider.Options, p)
		} else {
			genProvider, err = ToProvider(cnfProvider.Id, cnfProvider.Options)
		}

		if err != nil {
			return nstore, errors.Wrapf(err, "config to provider %s", cnfProvider.Name)
		}

		err = nstore.Store(cnfProvider.Name, genProvider)
		if err != nil {
			return nstore, errors.Wrap(err, "registering provider")
		}
	}
	return nstore, nil

}
