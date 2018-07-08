package provider

import (
	"time"
)

type Lease interface {
	End() Time
	HasResponse() bool
	RenewBefore() Time
	Response() (*Response, error)
}
