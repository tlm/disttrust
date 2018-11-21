package vault

import (
	"github.com/hashicorp/vault/api"
)

type Writer interface {
	Write(string, map[string]interface{}) (*api.Secret, error)
}

type WriterFunc func(string, map[string]interface{}) (*api.Secret, error)

func (w WriterFunc) Write(p string, d map[string]interface{}) (*api.Secret, error) {
	return w(p, d)
}
