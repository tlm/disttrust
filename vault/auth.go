package vault

import (
	"github.com/hashicorp/vault/api"
)

type AuthHandler interface {
	Auth(*api.Client) (*api.Secret, error)
}

type MakeAuthHandler func(map[string]string) (AuthHandler, error)

var (
	AuthHandlers = make(map[string]MakeAuthHandler)
)
