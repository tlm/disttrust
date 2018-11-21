package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type Config struct {
	Address string
	MakeCSR bool
	Path    string
	Role    string
	RollKey bool
}

type Provider struct {
	auth       AuthHandler
	authDoneCh chan error
	client     *api.Client
	config     Config
	issuer     Issuer
}

const (
	ProviderId = "vault"
)

func (p *Provider) authRenewal(secret *api.Secret) {
	renewer, _ := p.client.NewRenewer(&api.RenewerInput{
		Secret: secret,
	})
	go renewer.Renew()
	defer renewer.Stop()

	// At this stage we have setup a vault renewer that will keep renewing the
	// auth token for as long as possible. If the renewer tells us its done then
	// that means an error has occurred with renewal at which point this is
	// considered a failure or a new token needs to be generated. If a new token
	// generation also fails then nothing more we can do.
	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				p.authDoneCh <- errors.Wrap(err, "auth renewal")
				break
			}
			secret, err := p.auth.Auth(p.client)
			if err != nil {
				p.authDoneCh <- errors.Wrap(err, "making new auth token in auth renewal")
				break
			}
			p.client.SetToken(secret.Auth.ClientToken)
			go p.authRenewal(secret)
			break
		case renewal := <-renewer.RenewCh():
			p.client.SetToken(renewal.Secret.Auth.ClientToken)
		}
	}
}

func (p *Provider) Issue(req *provider.Request) (provider.Lease, error) {
	select {
	case err := <-p.authDoneCh:
		if err != nil {
			return nil, err
		}
	default:
	}

	return p.issuer.Issue(req)
}

func NewProvider(config Config, auth AuthHandler) (*Provider, error) {
	vconfig := api.DefaultConfig()
	vconfig.Address = config.Address

	client, err := api.NewClient(vconfig)
	if err != nil {
		return nil, errors.Wrap(err, "vault provider creation")
	}

	issuer := GenerateIssuer(config.Path, config.Role, client.Logical())

	nprv := Provider{
		auth:       auth,
		authDoneCh: make(chan error),
		client:     client,
		config:     config,
		issuer:     issuer,
	}

	tknSecret, err := nprv.auth.Auth(client)
	if err != nil {
		return nil, errors.Wrap(err, "new vault provider auth")
	}
	nprv.client.SetToken(tknSecret.Auth.ClientToken)
	go nprv.authRenewal(tknSecret)

	return &nprv, nil
}

func (p *Provider) Renew(lease provider.Lease) (provider.Lease, error) {
	select {
	case err := <-p.authDoneCh:
		if err != nil {
			return nil, err
		}
	default:
	}

	var vlease *Lease
	var ok bool
	if vlease, ok = lease.(*Lease); !ok {
		return nil, fmt.Errorf("unsupported lease type")
	}
	if vlease.renewable {
		// TODO implement
		return nil, fmt.Errorf("currently do not support renewable leases")
	}

	return p.Issue(vlease.Request())
}
