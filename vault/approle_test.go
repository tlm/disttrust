package vault

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func AuthMissingOptTest(t *testing.T) {
	tests := []struct {
		Opts map[string]string
	}{
		{map[string]string{"secretId": "123"}},
		{map[string]string{"roleId": "123"}},
	}

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatalf("creating dummy vault client %v", err)
	}

	for _, test := range tests {
		handle := AppRoleAuthHandler{}
		err := handle.Auth(client, test.Opts)
		if err == nil {
			t.Fatal("should have recieved error for missing opts")
		}
	}
}
