package config

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/spf13/viper"

	"github.com/tlmiller/disttrust/provider"
)

type providerConf struct {
	Id      string
	Name    string
	Options map[string]interface{}
}

const (
	Providers = "providers"
)

func GetProviderStore(v *viper.Viper,
	s *provider.Store) (*provider.Store, error) {
	return GetProviderStoreWithFactory(v, s, DefaultProviderFactory)
}

func GetProviderStoreWithFactory(v *viper.Viper, store *provider.Store,
	factory ProviderFactory) (*provider.Store, error) {
	nstore := provider.NewStore()
	providers := []providerConf{}
	err := v.UnmarshalKey(Providers, &providers)
	if err != nil {
		return nstore, errors.Wrap(err, "decoding providers config")
	}

	for i, p := range providers {
		if p.Name == "" {
			return nstore, fmt.Errorf("undefined provider name for index %d", i)
		}

		mapper, exists := factory[provider.Id(p.Id)]
		if !exists {
			return nil, fmt.Errorf("provider mapper does not exist for id %s", p.Id)
		}

		var genP provider.Provider
		changed := true

		// We are trying to work out if a previous provider exists with the same
		// config. If so we reuse it for reload support so we don't run auth
		// handlers again and potentially hit a token limit in a provider
		if prevP, exists := store.Fetch(p.Name); exists &&
			!HasProviderChanged(p.Options, prevP) {
			changed = false
			genP = prevP
		}

		if changed {
			var err error
			genP, err = mapper.MakeProvider(p.Options)
			if err != nil {
				return nstore, errors.Wrapf(err, "getting provider for %s", p.Name)
			}

			genP = &ProviderConfigWrap{
				Options: p.Options,
				P:       genP,
			}
		}

		err = nstore.Store(p.Name, genP)
		if err != nil {
			return nstore, errors.Wrapf(err, "registering provider %s", p.Name)
		}
	}
	return nstore, nil
}

func HasProviderChanged(opts map[string]interface{}, p provider.Provider) bool {
	confP, ok := p.(*ProviderConfigWrap)
	if !ok {
		return true
	}
	//Invert the deep equals because the provider has not changed
	return !reflect.DeepEqual(opts, confP.Options)
}
