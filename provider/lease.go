package provider

import (
	"time"
)

type Lease interface {
	ID() string
	HasResponse() bool
	Response() (*Response, error)
	Start() time.Time
	Till() time.Time
}
