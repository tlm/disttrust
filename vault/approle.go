package vault

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hashicorp/vault/api"
)

type AppRoleAuthHandler struct {
	roleId   string
	secretId string
}

func (a *AppRoleAuthHandler) Auth(client *api.Client) (*api.Secret, error) {
	data := map[string]interface{}{
		"role_id":   a.roleId,
		"secret_id": a.secretId,
	}

	path := "auth/approle/login/"
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return nil, errors.Wrap(err, "loging in with approle")
	}

	return secret, nil
}

func NewAppRoleAuthHandler(opts map[string]string) (*AppRoleAuthHandler, error) {
	roleId, exists := opts["roleId"]
	if !exists {
		return nil, fmt.Errorf("roleId not provided for approle auth handler")
	}
	secretId, exists := opts["secretId"]
	if !exists {
		return nil, fmt.Errorf("secretId not provided for approle auth handler")
	}
	return &AppRoleAuthHandler{
		roleId:   roleId,
		secretId: secretId,
	}, nil
}

func init() {
	AuthHandlers["approle"] =
		MakeAuthHandler(func(opts map[string]string) (AuthHandler, error) {
			return NewAppRoleAuthHandler(opts)
		})
}
