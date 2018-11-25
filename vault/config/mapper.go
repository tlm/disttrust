package config

import (
	"fmt"

	"github.com/pkg/errors"

	_ "github.com/spf13/viper"

	"github.com/tlmiller/disttrust/provider"
	"github.com/tlmiller/disttrust/request"
	"github.com/tlmiller/disttrust/util"
	"github.com/tlmiller/disttrust/vault"
)

type RequestConfig struct {
	CSR        bool
	KeyType    string
	RSABits    int
	ECDSACurve string
}

type VaultConfig struct {
	Address    string
	AuthMethod string
	AuthOpts   map[string]string
	Path       string
	Request    RequestConfig
	Role       string
}

const (
	RSAKeyType = "rsa"
)

var (
	DefaultRequestConfig = RequestConfig{
		CSR:        false,
		KeyType:    "rsa",
		RSABits:    2048,
		ECDSACurve: "p246",
	}
)

var (
	SupportedCurves   = []string{"p224", "p256", "p384", "p521"}
	SupportedKeyTypes = []string{"rsa", "ecdsa"}
)

func New() interface{} {
	return &VaultConfig{
		Request: DefaultRequestConfig,
	}
}

func Mapper(v interface{}) (provider.Provider, error) {
	conf, ok := v.(*VaultConfig)
	if !ok {
		return nil, errors.New("parsing vault provider config")
	}

	if conf.Path == "" {
		return nil, errors.New("no vault pki path specificed")
	}

	if conf.Role == "" {
		return nil, errors.New("no vault pki role for path specified")
	}

	var issuer vault.Issuer
	if conf.Request.CSR {
		if !util.StringInSlice(conf.Request.KeyType, SupportedKeyTypes) {
			return nil, fmt.Errorf(
				"request key type must be one of %v for vault csr requests",
				SupportedKeyTypes)
		}

		var keyMaker request.KeyMaker
		if conf.Request.KeyType == RSAKeyType {
			keyMaker = request.NewRSAKeyMaker(conf.Request.RSABits)
		}

		requester := request.NewCSRRequester(keyMaker)
		issuer = vault.CSRVerbatimIssuer(conf.Path, conf.Role, requester)
	} else {
		issuer = vault.GenerateIssuer(conf.Path, conf.Role)
	}

	authMaker, exists := vault.AuthHandlers[conf.AuthMethod]
	if !exists {
		return nil, fmt.Errorf("no auth handler for method '%s'", conf.AuthMethod)
	}
	auth, err := authMaker(conf.AuthOpts)
	if err != nil {
		return nil, errors.Wrap(err, "making auth handler")
	}

	pconfig := vault.Config{}
	pconfig.Address = conf.Address

	provider, err := vault.NewProvider(pconfig, issuer, auth)
	if err != nil {
		return nil, errors.Wrap(err, "building vault provider from config")
	}
	return provider, nil
}
