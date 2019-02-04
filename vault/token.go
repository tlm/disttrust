package vault

import (
	"github.com/hashicorp/vault/api"

	"github.com/pkg/errors"
)

func TokenValid(client *api.Client, tknSecret *api.Secret) (bool, error) {
	client.SetToken(tknSecret.Auth.ClientToken)

	_, err := client.Logical().ReadWithData("/auth/token/lookup-self", nil)
	if err != nil {
		return false, errors.Wrap(err, "testing token validity")
	}

	return true, nil
}
