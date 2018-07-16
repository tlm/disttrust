package vault

import (
	"testing"
)

func TestAuthMissingOpt(t *testing.T) {
	tests := []struct {
		Opts map[string]string
	}{
		{map[string]string{"secretId": "123"}},
		{map[string]string{"roleId": "123"}},
	}

	for _, test := range tests {
		_, err := NewAppRoleAuthHandler(test.Opts)
		if err == nil {
			t.Fatal("should have recieved error for missing opts")
		}
	}
}
