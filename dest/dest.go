package dest

import (
	"github.com/tlmiller/disttrust/provider"
)

type Dest interface {
	Send(*provider.Response) error
}
