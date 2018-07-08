package provider

import (
	"time"
)

type Lease interface {
	HasResponse() bool
	Response() (*Response, error)
	Till() time.Time
}
