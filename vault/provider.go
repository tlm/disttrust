package vault

import (
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type AuthHandler interface {
	Auth(*api.Client, map[string]string) error
}

type Config struct {
	Address string
	Path    string
	Role    string
}

type Provider struct {
	client *api.Client
	config Config
}

const (
	ProviderId = "vault"
)

var (
	AuthHandlers = make(map[string]AuthHandler)
)

func DefaultConfig() *api.Config {
	return api.DefaultConfig()
}

func (p *Provider) Issue(req *provider.Request) (provider.Lease, error) {
	logical := p.client.Logical()
	path := fmt.Sprintf("%s/issue/%s", p.config.Path, p.config.Role)

	data := map[string]interface{}{}
	data["alt_names"] = strings.Join(req.AltNames, ",")
	data["common_name"] = req.CommonName
	data["format"] = "pem"

	secret, err := logical.Write(path, data)
	if err != nil {
		return nil, errors.Wrap(err, "issue certificate from vault")
	}

	lease, err := LeaseFromSecret(secret)
	if err != nil {
		return nil, errors.Wrap(err, "issue certificate make lease")
	}
	return lease, nil
}

func NewProvider(config Config, auth AuthHandler,
	authOpt map[string]string) (*Provider, error) {

	vconfig := api.DefaultConfig()
	vconfig.Address = config.Address

	client, err := api.NewClient(vconfig)
	if err != nil {
		return nil, errors.Wrap(err, "vault provider creation")
	}

	err = auth.Auth(client, authOpt)
	if err != nil {
		return nil, errors.Wrap(err, "new vault provider auth")
	}

	return &Provider{
		client: client,
		config: config,
	}, nil
}

func Renew(lease provider.Lease) (Lease, error) {
	return nil, errors.New("not  implemented")
}
