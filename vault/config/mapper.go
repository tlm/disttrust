package config

import (
	"fmt"

	"github.com/pkg/errors"

	_ "github.com/spf13/viper"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/vault"
)

type VaultConfig struct {
	Address    string
	AuthMethod string
	AuthOpts   map[string]string
	Path       string
	Role       string
}

func New() interface{} {
	return &VaultConfig{}
}

func Mapper(v interface{}) (provider.Provider, error) {
	conf, ok := v.(*VaultConfig)
	if !ok {
		return nil, errors.New("parsing vault provider config")
	}

	authMaker, exists := vault.AuthHandlers[conf.AuthMethod]
	if !exists {
		return nil, fmt.Errorf("no auth handler for method '%s'", conf.AuthMethod)
	}
	auth, err := authMaker(conf.AuthOpts)
	if err != nil {
		return nil, errors.Wrap(err, "making auth handler")
	}

	pconfig := vault.Config{}
	pconfig.Address = conf.Address

	if conf.Path == "" {
		return nil, errors.New("no vault pki path specificed")
	}
	pconfig.Path = conf.Path
	if conf.Role == "" {
		return nil, errors.New("no vault pki role for path specified")
	}
	pconfig.Role = conf.Role

	provider, err := vault.NewProvider(pconfig, auth)
	if err != nil {
		return nil, errors.Wrap(err, "building vault provider from config")
	}
	return provider, nil
}
