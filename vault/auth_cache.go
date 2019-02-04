package vault

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hashicorp/vault/api"

	"github.com/pkg/errors"

	"github.com/tlmiller/disttrust/file"
)

type AuthCache interface {
	Read() (*api.Secret, error)
	Write(*api.Secret) error
}

type FileAuthCache struct {
	secretFile file.File
}

type EmptyAuthCache struct {
}

func NewFileAuthCache(secretFile file.File) *FileAuthCache {
	return &FileAuthCache{
		secretFile: secretFile,
	}
}

func (e *EmptyAuthCache) Read() (*api.Secret, error) {
	return nil, nil
}

func (f *FileAuthCache) Read() (*api.Secret, error) {
	buf, err := ioutil.ReadFile(f.secretFile.Path)
	if err != nil && os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	secret := &api.Secret{}
	if err := json.Unmarshal(buf, secret); err != nil {
		return secret, err
	}
	return secret, err
}

func (e *EmptyAuthCache) Write(_ *api.Secret) error {
	return nil
}

func (f *FileAuthCache) Write(secret *api.Secret) error {
	buf, err := json.Marshal(secret)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(f.secretFile.Path, buf, f.secretFile.Mode); err != nil {
		return err
	}
	if err := f.secretFile.Chown(); err != nil {
		return errors.Wrapf(err, "chowning secret file %s", f.secretFile.Path)
	}
	return nil
}
