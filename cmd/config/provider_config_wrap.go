package config

import (
	"encoding/json"

	"github.com/tlmiller/disttrust/provider"
)

type ProviderConfigWrap struct {
	config json.RawMessage
	p      provider.Provider
}

func (p *ProviderConfigWrap) Issue(r *provider.Request) (provider.Lease, error) {
	return p.p.Issue(r)
}

func (p *ProviderConfigWrap) Renew(l provider.Lease) (provider.Lease, error) {
	return p.p.Renew(l)
}
