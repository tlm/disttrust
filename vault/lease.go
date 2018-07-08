package vault

import (
	"time"

	"github.com/hashicorp/vault/api"

	"github.com/tlmiller/disttrust/provider"
)

type Lease struct {
	end       Time
	leaseId   string
	renewable bool
	renewBy   Time
	response  *provider.Response
}

func NewLeaseFromSecret(secret *api.Secret) (*Lease, error) {
}
