package request

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
)

type rsaKey struct {
	key *rsa.PrivateKey
}

type RSAKeyMaker struct {
	Bits int
}

func (r *RSAKeyMaker) MakeKey() (Key, error) {
	k, err := rsa.GenerateKey(rand.Reader, r.Bits)
	return &rsaKey{
		key: k,
	}, err
}

func NewRSAKeyMaker(bits int) *RSAKeyMaker {
	return &RSAKeyMaker{
		Bits: bits,
	}
}

func (k *rsaKey) PKCS8() ([]byte, error) {
	return x509.MarshalPKCS8PrivateKey(k.key)
}

func (k *rsaKey) Raw() interface{} {
	return k.key
}
