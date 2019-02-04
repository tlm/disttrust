package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/provider"
)

type Provider struct {
	auth       AuthHandler
	authCache  AuthCache
	authDoneCh chan error
	client     *api.Client
	issuer     Issuer
	name       string
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
finish:
	for {
		select {
		case err := <-renewer.DoneCh():
			if err != nil {
				p.authDoneCh <- errors.Wrap(err, "auth renewal")
				close(p.authDoneCh)
				break finish
			}
			secret, err := p.auth.Auth(p.client)
			if err != nil {
				p.authDoneCh <- errors.Wrap(err, "making new auth token in auth renewal")
				close(p.authDoneCh)
				break finish
			}

			if err := p.authCache.Write(secret); err != nil {
				p.authDoneCh <- errors.Wrap(err, "caching auth secret")
				close(p.authDoneCh)
				break finish
			}
			p.client.SetToken(secret.Auth.ClientToken)
			go p.authRenewal(secret)
			break finish
		case renewal := <-renewer.RenewCh():
			if err := p.authCache.Write(renewal.Secret); err != nil {
				p.authDoneCh <- errors.Wrap(err, "caching auth secret")
				close(p.authDoneCh)
				break finish
			}
			p.client.SetToken(renewal.Secret.Auth.ClientToken)
		}
	}
}

func (p *Provider) Initialise() error {
	tknSecret, err := p.authCache.Read()
	if err != nil {
		return errors.Wrap(err, "getting auth cache secret for init")
	}

	if tknSecret != nil {
		_, err := TokenValid(p.client, tknSecret)
		if err != nil {
			tknSecret = nil
		}
	}

	if tknSecret == nil {
		tknSecret, err = p.auth.Auth(p.client)
		if err != nil {
			return errors.Wrap(err, "new vault provider auth")
		}
		if err := p.authCache.Write(tknSecret); err != nil {
			return errors.Wrap(err, "caching auth secret for new provider")
		}
	}

	p.client.SetToken(tknSecret.Auth.ClientToken)
	go p.authRenewal(tknSecret)
	return nil
}

func (p *Provider) Issue(req *provider.Request) (provider.Lease, error) {
	select {
	case err, ok := <-p.authDoneCh:
		if err != nil || !ok {
			return nil, err
		}
	default:
	}

	return p.issuer.Issue(req, p.client.Logical())
}

func (p *Provider) Name() string {
	return p.name
}

func NewProvider(name, address string, issuer Issuer, auth AuthHandler,
	cache AuthCache) (*Provider, error) {
	vconfig := api.DefaultConfig()
	vconfig.Address = address

	client, err := api.NewClient(vconfig)
	if err != nil {
		return nil, errors.Wrap(err, "vault provider creation")
	}

	nprv := Provider{
		auth:       auth,
		authCache:  cache,
		authDoneCh: make(chan error),
		client:     client,
		issuer:     issuer,
		name:       name,
	}

	return &nprv, nil
}

func (p *Provider) Renew(lease provider.Lease) (provider.Lease, error) {
	select {
	case err, ok := <-p.authDoneCh:
		if err != nil || !ok {
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
