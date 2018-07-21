package config

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/vault"
)

type vaultConfig struct {
	Address    string            `json:"address"`
	AuthMethod string            `json:"authMethod"`
	AuthOpts   map[string]string `json:"authOpts"`
	Path       string            `json:"path"`
	Role       string            `json:"role"`
}

func init() {
	err := MapProvider(vault.ProviderId, ProviderMapper(vaultMapper))
	if err != nil {
		panic(err)
	}
}

func vaultMapper(opt json.RawMessage) (provider.Provider, error) {
	config := vaultConfig{}
	err := json.Unmarshal(opt, &config)
	if err != nil {
		return nil, errors.Wrap(err, "parsing vault provider config")
	}

	authMaker, exists := vault.AuthHandlers[config.AuthMethod]
	if !exists {
		return nil, fmt.Errorf("no auth handler for method '%s'", config.AuthMethod)
	}
	auth, err := authMaker(config.AuthOpts)
	if err != nil {
		return nil, errors.Wrap(err, "making auth handler")
	}

	pconfig := vault.Config{}
	pconfig.Address = config.Address

	if len(config.Path) == 0 {
		return nil, errors.New("no vault pki path specificed")
	}
	pconfig.Path = config.Path
	if len(config.Role) == 0 {
		return nil, errors.New("no vault pki role for path specified")
	}
	pconfig.Role = config.Role

	provider, err := vault.NewProvider(pconfig, auth)
	if err != nil {
		return nil, errors.Wrap(err, "building vault provider from config")
	}
	return provider, nil
}
