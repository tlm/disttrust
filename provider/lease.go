package provider

import (
	"time"
)

type Lease interface {
	HasResponse() bool
	Response() (*Response, error)
	Start() time.Time
	Till() time.Time
}
