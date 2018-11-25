package request

import (
	"testing"
)

func TestRSAKeyGenMeta(t *testing.T) {
	maker := NewRSAKeyMaker(4096)
	_, err := maker.MakeKey()
	if err != nil {
		t.Fatalf("unexpected error making new rsa key: %v", err)
	}
}
