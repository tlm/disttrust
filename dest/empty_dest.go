package dest

import (
	"github.com/tlmiller/disttrust/provider"
)

type Empty struct {
}

func (e *Empty) Send(_ *provider.Response) error {
	return nil
}
