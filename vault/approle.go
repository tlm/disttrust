package vault

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hashicorp/vault/api"
)

type AppRoleAuthHandler struct {
}

func (a *AppRoleAuthHandler) Auth(client *api.Client, opt map[string]string) error {
	roleId, exists := opt["roleId"]
	if !exists {
		return fmt.Errorf("roleId not provided for approle auth handler")
	}
	secretId, exists := opt["secretId"]
	if !exists {
		return fmt.Errorf("secretId not provided for approle auth handler")
	}

	data := map[string]interface{}{
		"role_id":   roleId,
		"secret_id": secretId,
	}

	path := "auth/approle/login/"
	secret, err := client.Logical().Write(path, data)
	if err != nil {
		return errors.Wrap(err, "loging in with approle")
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}

func init() {
	AuthHandlers["approle"] = &AppRoleAuthHandler{}
}
